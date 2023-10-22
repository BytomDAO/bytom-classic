Anonimitycash version 1.0.10 is now available from:

  https://github.com/Anonimitycash/anonimitycash/releases/tag/v1.0.10


Please report bugs using the issue tracker at github:

  https://github.com/Anonimitycash/anonimitycash/issues

How to Upgrade
===============

If you are running an older version, shut it down. Wait until it has quited completely, and then run the new version Anonimitycash.
You can operate according to the user manual.[(Anonimitycash User Manual)](https://anonimitycash.io/wp-content/themes/freddo/images/wallet/AnonimitycashUsermanualV1.0_en.pdf)


1.0.10 changelog
================
__Anonimitycash Node__

+ [`PR #1738`](https://github.com/Anonimitycash/anonimitycash/pull/1738)
    - Add the core block intergra testing case. 
        - Including data correction in LevelDB, memory and orphan after block processing.
+ [`PR #1745`](https://github.com/Anonimitycash/anonimitycash/pull/1745) 
    - Add the core block intergra testing case. 
        - Including attach a block, 
        - process an orphan block, 
        - adding an block into forked chain,
        - adding an block casing rollback
        - and all other combines of transactions type in block.
+ [`PR #1751`](https://github.com/Anonimitycash/anonimitycash/pull/1751)
    - Fixed synced error between node block mining. 
+ [`PR #1777`](https://github.com/Anonimitycash/anonimitycash/pull/1777)
    - Fixed the transactions failed re-entered the transactions pool when chain reorganized.
+ [`PR #1780`](https://github.com/Anonimitycash/anonimitycash/pull/1780) 
    - Fixed the banned node forbidden error.
+ [`PR #1789`](https://github.com/Anonimitycash/anonimitycash/pull/1789)
    - Add the handshake permision for Ed25519 node only.
+ [`PR #1791`](https://github.com/Anonimitycash/anonimitycash/pull/1791)
    - add `/estimate-chain-tx-gas` API to estimate chain transactions gas when building chain transactions 
+ [`PR #1792`](https://github.com/Anonimitycash/anonimitycash/pull/1792) 
    - fix `/estimate-chain-tx-gas` API response format inconsistency in case the (chain) transaction to build contains only ONE transaction

__Anonimitycash Dashboard__

+ [`PR #1798`](https://github.com/Anonimitycash/anonimitycash/pull/1798) 
    - Update dashboard with estimate chain transactions fee function and a switcher for either the chain transactions or the normal transactions. This feature supports by offical MITY asset only.

Credits
--------

Thanks to everyone who directly contributed to this release:

- Agouri
- Colt-Z
- HAOYUatHZ
- langyu
- Paladz
- shenao78
- shengling2008
- yahtoo
- zcc0721
- ZhitingLin

And everyone who helped test.
