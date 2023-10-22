package node

import (
	"context"
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"

	"github.com/prometheus/prometheus/util/flock"
	log "github.com/sirupsen/logrus"
	cmn "github.com/tendermint/tmlibs/common"
	browser "github.com/toqueteos/webbrowser"

	"github.com/anonimitycash/anonimitycash-classic/accesstoken"
	"github.com/anonimitycash/anonimitycash-classic/account"
	"github.com/anonimitycash/anonimitycash-classic/api"
	"github.com/anonimitycash/anonimitycash-classic/asset"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/pseudohsm"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/txfeed"
	cfg "github.com/anonimitycash/anonimitycash-classic/config"
	"github.com/anonimitycash/anonimitycash-classic/consensus"
	"github.com/anonimitycash/anonimitycash-classic/database"
	dbm "github.com/anonimitycash/anonimitycash-classic/database/leveldb"
	"github.com/anonimitycash/anonimitycash-classic/env"
	"github.com/anonimitycash/anonimitycash-classic/event"
	anonimitycashLog "github.com/anonimitycash/anonimitycash-classic/log"
	"github.com/anonimitycash/anonimitycash-classic/mining/cpuminer"
	"github.com/anonimitycash/anonimitycash-classic/mining/miningpool"
	"github.com/anonimitycash/anonimitycash-classic/mining/tensority"
	"github.com/anonimitycash/anonimitycash-classic/net/websocket"
	"github.com/anonimitycash/anonimitycash-classic/netsync"
	"github.com/anonimitycash/anonimitycash-classic/p2p"
	"github.com/anonimitycash/anonimitycash-classic/protocol"
	w "github.com/anonimitycash/anonimitycash-classic/wallet"
)

const (
	webHost   = "http://127.0.0.1"
	logModule = "node"
)

// Node represent anonimitycash node
type Node struct {
	cmn.BaseService

	config          *cfg.Config
	eventDispatcher *event.Dispatcher
	syncManager     *netsync.SyncManager

	wallet          *w.Wallet
	accessTokens    *accesstoken.CredentialStore
	notificationMgr *websocket.WSNotificationManager
	api             *api.API
	chain           *protocol.Chain
	txfeed          *txfeed.Tracker
	cpuMiner        *cpuminer.CPUMiner
	miningPool      *miningpool.MiningPool
	miningEnable    bool
}

// NewNode create anonimitycash node
func NewNode(config *cfg.Config) *Node {
	ctx := context.Background()
	if err := lockDataDirectory(config); err != nil {
		cmn.Exit("Error: " + err.Error())
	}

	if err := anonimitycashLog.InitLogFile(config); err != nil {
		log.WithField("err", err).Fatalln("InitLogFile failed")
	}

	initActiveNetParams(config)
	initCommonConfig(config)

	// Get store
	if config.DBBackend != "memdb" && config.DBBackend != "leveldb" {
		cmn.Exit(cmn.Fmt("Param db_backend [%v] is invalid, use leveldb or memdb", config.DBBackend))
	}
	coreDB := dbm.NewDB("core", config.DBBackend, config.DBDir())
	store := database.NewStore(coreDB)

	tokenDB := dbm.NewDB("accesstoken", config.DBBackend, config.DBDir())
	accessTokens := accesstoken.NewStore(tokenDB)

	dispatcher := event.NewDispatcher()
	txPool := protocol.NewTxPool(store, dispatcher)
	chain, err := protocol.NewChain(store, txPool)
	if err != nil {
		cmn.Exit(cmn.Fmt("Failed to create chain structure: %v", err))
	}

	var accounts *account.Manager
	var assets *asset.Registry
	var wallet *w.Wallet
	var txFeed *txfeed.Tracker

	txFeedDB := dbm.NewDB("txfeeds", config.DBBackend, config.DBDir())
	txFeed = txfeed.NewTracker(txFeedDB, chain)

	if err = txFeed.Prepare(ctx); err != nil {
		log.WithFields(log.Fields{"module": logModule, "error": err}).Error("start txfeed")
		return nil
	}

	hsm, err := pseudohsm.New(config.KeysDir())
	if err != nil {
		cmn.Exit(cmn.Fmt("initialize HSM failed: %v", err))
	}

	if !config.Wallet.Disable {
		walletDB := dbm.NewDB("wallet", config.DBBackend, config.DBDir())
		accounts = account.NewManager(walletDB, chain)
		assets = asset.NewRegistry(walletDB, chain)
		wallet, err = w.NewWallet(walletDB, accounts, assets, hsm, chain, dispatcher, config.Wallet.TxIndex)
		if err != nil {
			log.WithFields(log.Fields{"module": logModule, "error": err}).Error("init NewWallet")
		}

		// trigger rescan wallet
		if config.Wallet.Rescan {
			wallet.RescanBlocks()
		}
	}

	syncManager, err := netsync.NewSyncManager(config, chain, txPool, dispatcher)
	if err != nil {
		cmn.Exit(cmn.Fmt("Failed to create sync manager: %v", err))
	}

	notificationMgr := websocket.NewWsNotificationManager(config.Websocket.MaxNumWebsockets, config.Websocket.MaxNumConcurrentReqs, chain, dispatcher)

	// run the profile server
	profileHost := config.ProfListenAddress
	if profileHost != "" {
		// Profiling anonimitycashd programs.see (https://blog.golang.org/profiling-go-programs)
		// go tool pprof http://profileHose/debug/pprof/heap
		go func() {
			if err = http.ListenAndServe(profileHost, nil); err != nil {
				cmn.Exit(cmn.Fmt("Failed to register tcp profileHost: %v", err))
			}
		}()
	}

	node := &Node{
		eventDispatcher: dispatcher,
		config:          config,
		syncManager:     syncManager,
		accessTokens:    accessTokens,
		wallet:          wallet,
		chain:           chain,
		txfeed:          txFeed,
		miningEnable:    config.Mining,

		notificationMgr: notificationMgr,
	}

	node.cpuMiner = cpuminer.NewCPUMiner(chain, accounts, txPool, dispatcher)
	node.miningPool = miningpool.NewMiningPool(chain, accounts, txPool, dispatcher)

	node.BaseService = *cmn.NewBaseService(nil, "Node", node)

	if config.Simd.Enable {
		tensority.UseSIMD = true
	}

	return node
}

// Lock data directory after daemonization
func lockDataDirectory(config *cfg.Config) error {
	_, _, err := flock.New(filepath.Join(config.RootDir, "LOCK"))
	if err != nil {
		return errors.New("datadir already used by another process")
	}
	return nil
}

func initActiveNetParams(config *cfg.Config) {
	var exist bool
	consensus.ActiveNetParams, exist = consensus.NetParams[config.ChainID]
	if !exist {
		cmn.Exit(cmn.Fmt("chain_id[%v] don't exist", config.ChainID))
	}
}

func initCommonConfig(config *cfg.Config) {
	cfg.CommonConfig = config
}

// Lanch web broser or not
func launchWebBrowser(port string) {
	webAddress := webHost + ":" + port
	log.Info("Launching System Browser with :", webAddress)
	if err := browser.Open(webAddress); err != nil {
		log.Error(err.Error())
		return
	}
}

func (n *Node) initAndstartAPIServer() {
	n.api = api.NewAPI(n.syncManager, n.wallet, n.txfeed, n.cpuMiner, n.miningPool, n.chain, n.config, n.accessTokens, n.eventDispatcher, n.notificationMgr)

	listenAddr := env.String("LISTEN", n.config.ApiAddress)
	env.Parse()
	n.api.StartServer(*listenAddr)
}

func (n *Node) OnStart() error {
	if n.miningEnable {
		if _, err := n.wallet.AccountMgr.GetMiningAddress(); err != nil {
			n.miningEnable = false
			log.Error(err)
		} else {
			n.cpuMiner.Start()
		}
	}
	if !n.config.VaultMode {
		if err := n.syncManager.Start(); err != nil {
			return err
		}
	}

	n.initAndstartAPIServer()
	if err := n.notificationMgr.Start(); err != nil {
		return err
	}

	if !n.config.Web.Closed {
		_, port, err := net.SplitHostPort(n.config.ApiAddress)
		if err != nil {
			log.Error("Invalid api address")
			return err
		}
		launchWebBrowser(port)
	}
	return nil
}

func (n *Node) OnStop() {
	n.notificationMgr.Shutdown()
	n.notificationMgr.WaitForShutdown()
	n.BaseService.OnStop()
	if n.miningEnable {
		n.cpuMiner.Stop()
	}
	if !n.config.VaultMode {
		n.syncManager.Stop()
	}
	n.eventDispatcher.Stop()
}

func (n *Node) RunForever() {
	// Sleep forever and then...
	cmn.TrapSignal(func() {
		n.Stop()
	})
}

func (n *Node) NodeInfo() *p2p.NodeInfo {
	return n.syncManager.NodeInfo()
}

func (n *Node) MiningPool() *miningpool.MiningPool {
	return n.miningPool
}
