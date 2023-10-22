Anonimitycash version 1.1.0 is now available from:

  https://github.com/Anonimitycash/anonimitycash/releases/tag/v1.1.0


Please report bugs using the issue tracker at github:

  https://github.com/Anonimitycash/anonimitycash/issues

How to Upgrade
===============

Please notice new version anonimitycash path is $GOPATH/src/github.com/anonimitycash/anonimitycash if you install anonimitycash from source code.  
If you are running an older version, shut it down. Wait until it has quited completely, and then run the new version Anonimitycash.
You can operate according to the user manual.[(Anonimitycash User Manual)](https://anonimitycash.io/wp-content/themes/freddo/images/wallet/AnonimitycashUsermanualV1.0_en.pdf)


1.1.0 changelog
================
__Anonimitycash Node__

+ [`PR #1805`](https://github.com/Anonimitycash/anonimitycash/pull/1805)
    - Correct anonimitycash go import path to github.com/anonimitycash/anonimitycash. Developer can use go module to manage dependency of anonimitycash. 
+ [`PR #1815`](https://github.com/Anonimitycash/anonimitycash/pull/1815) 
    - Add asynchronous validate transactions function to optimize the performance of validating and saving block. 

__Anonimitycash Dashboard__

+ [`PR #1829`](https://github.com/Anonimitycash/anonimitycash/pull/1829) 
    - Fixed the decimals type string to integer in create asset page.

Credits
--------

Thanks to everyone who directly contributed to this release:

- DeKaiju
- iczc
- Paladz
- zcc0721
- ZhitingLin

And everyone who helped test.
