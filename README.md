# ![Horcrux Logo](logo.png)

Horcrux is a command line tool that uses [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) to split a file into `n` encrypted parts, requiring `m` of those parts to later reconstruct. Horcrux is a usability layer built around Hashicorp Vault's [shamirs code](https://github.com/hashicorp/vault/tree/main/shamir) which is vendored [here](https://gitlab.com/unit410/vault-shamir).

## Split

A simple (non-encrypted) split is performed by specifying `--num-shares` and an input file. `--threshold` specifies how many shares are required to rebuild the whole and has a default value of `2`.

```shell
horcrux split --num-shares 5 --threshold 3 /path/to/secrets.txt
```

This will create 5 files in the current directory, any 3 of which can be combined to reconstruct the original plaintext. Each file contains the share payload and a threshold.

```shell
file.0.json
file.1.json
file.2.json
file.3.json
file.4.json
```

### Encrypted Split

You can optionally encrypt each share using a set of gpg public keys; one share per gpg key.

```shell
horcrux split --gpg-keys /folder/of/gpg/key/files/ --threshold 3 /path/to/secrets.txt
```

This will also create json files in the current directory, but each file will be encrypted to a unique key file from the `--gpg-keys` folder. Once decrypted, any 3 of these files can be used to reconstruct the original plaintext.

## Restore

Gather at least [threshold] number of files for the restore procedure and invoke `horcrux restore` providing the files as arguments.

```shell
horcrux restore file.0.json file.1.json file.4.json
```

If you have enough files (at least threshold) this will assemble the files back into the original and print the contents to screen. If the contents are not as expected - double check that you have enough share files to meet the threshold. If these files are encrypted, Horcrux will prompt to ensure private keys are loaded or smartcards are inserted before attempting to decrypt.
