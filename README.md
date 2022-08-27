# `ussher`

`ussher` aims to provide a backend for `sshd`'s [`AuthorizedKeysCommand`
option](https://man.openbsd.org/sshd_config.5#AuthorizedKeysCommand), by
remotely sourcing SSH `authorized_keys`.

> `AuthorizedKeysCommand`: Specifies a program to be used to look up the user's
  public keys. The program must be owned by root, not writable by group or
  others and specified by an absolute path. Arguments to
  `AuthorizedKeysCommand` accept the tokens described in the TOKENS section. If
  no arguments are specified then the username of the target user is used.
> The program should produce on standard output zero or more lines of
  authorized_keys output. `AuthorizedKeysCommand` is tried after the usual
  `AuthorizedKeysFile` files and will not be executed if a matching key is
  found there. By default, no `AuthorizedKeysCommand` is run.

## Configuration

`~/.ssh/authorized_keys.yml`:

```
urls:
- https://github.com/dolph.keys
```

## Usage

```
$ ussher
```

## Building

```
$ go build
```
