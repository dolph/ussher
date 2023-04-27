# `ussher`

`ussher` aims to provide a backend for `sshd`'s [`AuthorizedKeysCommand` option](https://man.openbsd.org/sshd_config.5#AuthorizedKeysCommand), by remotely sourcing SSH `authorized_keys`. In short:

> `AuthorizedKeysCommand`: Specifies a program to be used to look up the user's public keys. The program must be owned by root, not writable by group or others and specified by an absolute path.
>
> The program should produce on standard output zero or more lines of authorized_keys output. `AuthorizedKeysCommand` is tried after the usual `AuthorizedKeysFile` files and will not be executed if a matching key is found there. By default, no `AuthorizedKeysCommand` is run.

When `~/.ssh/authorized_keys` does not contain the keys required to authenticate a user, `sshd` invokes `ussher` to provide additional, remotely-sourced keys, such as from Github or another identity and access management provider.

## How it works

`ussher` provides a fallback mechanism for statically-defined `authorized_keys` files, such as when `authorized_keys` needs to be frequently updated, composed from a large number of sources, or simply defined at the moment of authorization.

1. When you `ssh $USER@$HOSTNAME`, `sshd` first reads something like `/home/.ssh/authorized_keys` to authenticate incoming SSH connections. If `authorized_keys` contains a public key that matches the incoming connection, then the connection attempt proceeds normally.

1. If that file does not exist or does not contain a public key for the incoming connection, `sshd` will invoke `AuthorizedKeysCommand` as `AuthorizedKeysCommandUser`. The `AuthorizedKeysCommand` (`ussher`, in this case), is responsible for returning a list of authorized public keys, which `sshd` then uses to validate the incoming connection.

1. When invoked, `ussher` sources authorized keys from any number of remote sources, such as a static text file or more commonly something like `https://github.com/{username}.keys`. `ussher` then returns a single set of authorized keys to `sshd` which are used to validate the incoming connection.

## Recommended installation & usage

1. Create a dedicated user and group to run `ussher`, named `ussher`:

   ```bash
   sudo adduser --system --user-group ussher
   ```

2. Install the `ussher` binary: to `/usr/local/bin`.

   ```bash
   sudo install -o root -g ussher -m 0750 ussher /usr/local/bin/ussher
   ```

3. Create a configuration directory.

   ```bash
   sudo mkdir --parents /etc/ussher
   ```

4. Create a directory for caching remotely-sourced data.

   ```bash
   sudo mkdir --parents /var/cache
   sudo mkdir --parents --mode=0700 /var/cache/ussher
   sudo chown ussher:ussher /var/cache/ussher
   ```

5. Create a directory for logging.

   ```bash
   sudo mkdir --parents --mode=0700 /var/log/ussher
   sudo chown ussher:ussher /var/log/ussher
   ```

6. Configure `sshd` to invoke `ussher`. Add the following lines to
   `/etc/ssh/sshd_config`:

   ```
   AuthorizedKeysCommand /usr/local/bin/ussher
   AuthorizedKeysCommandUser ussher
   ```

   You can script this with:

   ```bash
   sudo sed -i -E "s~^#?AuthorizedKeysCommand .*~AuthorizedKeysCommand /usr/local/bin/ussher~" /etc/ssh/sshd_config
   sudo sed -i -E "s~^#?AuthorizedKeysCommandUser .*~AuthorizedKeysCommandUser ussher~" /etc/ssh/sshd_config
   ```

   Validate sshd's new configuration and restart `sshd`, for example:

   ```bash
   sudo sshd -t
   sudo systemctl restart sshd
   ```

## Configuration

`/etc/ussher` contains configuration files for each user it supports. For example, to allow @dolph to SSH to your host as root (but, you know, _don't_), you would configure `/etc/ussher/root.yml` using:

```yaml
sources:
- url: https://github.com/dolph.keys
```

## Troubleshooting

### `Refusing to run unnecessarily writable binary`

Per the `sshd` man page:

> The program must be owned by root, not writable by group or others and specified by an absolute path.

`ussher` checks to ensure it's own binary is not unnecessarily writable at startup. If it is, a malicious user could remove this check or return any set of keys to sshd, which may be difficult to detect.

Ensure the file mode is similar to:

```
$ ls -l /usr/local/bin/ussher
-rwxr-x---. 1 root ussher 7823184 Apr 27 11:10 /usr/local/bin/ussher*
```

For example:

```bash
sudo chmod g-w /usr/local/bin/ussher
sudo chmod o-w /usr/local/bin/ussher
```

... but the binary may have been tainted. Verify it's checksum.

### `Refusing to run as root`

`ussher` checks to ensure it's not unnecessarily _running_ as root.

Ensure `sshd_config` specifies a non-root user:

```
AuthorizedKeysCommandUser ussher
```

### `Failed to write to /var/log/ussher`

`ussher` tries to ensure it can produce file-based logging output for auditing purposes. It's first preference is for `/var/log/ussher` which requires:

```bash
sudo mkdir --parents --mode=0700 /var/log/ussher
sudo chown ussher:ussher /var/log/ussher
```

### `Refusing to run without being able to log to /var/log/ussher/ or current working directory`

`ussher` tries to ensure it can produce file-based logging output for auditing purposes, and will fallback to the current working directory if `/var/log/ussher` is not writable. The best solution is to ensure `/var/log/ussher` is writable (see previous troubleshooting issue).

### `usage: ussher <username>`

Per the `sshd_config` man page:

> Arguments to AuthorizedKeysCommand accept the tokens described in the TOKENS section. If no arguments are specified then the username of the target user is used.

`ussher` only expects the default configuration from sshd (the username of the target user).

Ensure `sshd_config` only specifies the absolute path to `ussher`, without any additional arguments:

```
AuthorizedKeysCommand /usr/local/bin/ussher
```

### `User not found`

The user specified to `ussher` is either not a valid Linux username or not an existing user on the host. Double check the username specified to `ussher` as well as the `AuthorizedKeysCommand` value in `/etc/ssh/sshd_config`.
