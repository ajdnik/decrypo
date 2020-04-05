# decrypo

[![Build Status](https://github.com/ajdnik/decrypo/workflows/push-to-master/badge.svg "GitHub Actions status")](https://github.com/ajdnik/decrypo/actions?query=workflow%3Apush-to-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ajdnik/decrypo)](https://goreportcard.com/report/github.com/ajdnik/decrypo)

Decrypt Pluralsight videos into .mp4 format.

### Install

##### macOS with homebrew

First get the homebrew tap and install the formula

```bash
brew tap ajdnik/decrypo
brew install decrypo
```

### Usage

Run the command from your terminal application and define an `-output` flag where the decrypted videos should be stored.

```bash
$ decrypo -output "./Course Videos/"
Found 20 clips in database.
Decrypting clips and extracting captions...
20 / 20 [------------------------------------------------>] 100% 35 p/s
Successfully decrypted 20 of 20 clips.
You can find them at ./Course Videos/
```

To find out more about other flags use the `decrypo --help`.
