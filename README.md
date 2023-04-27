# `ussher`

`ussher` aims to provide a backend for `sshd`'s [`AuthorizedKeysCommand`
option](https://man.openbsd.org/sshd_config.5#AuthorizedKeysCommand), by
remotely sourcing SSH `authorized_keys`. In short:

> `AuthorizedKeysCommand`: Specifies a program to be used to look up the user's
  public keys. The program must be owned by root, not writable by group or
  others and specified by an absolute path.
>
> The program should produce on standard output zero or more lines of
  authorized_keys output. `AuthorizedKeysCommand` is tried after the usual
  `AuthorizedKeysFile` files and will not be executed if a matching key is
  found there. By default, no `AuthorizedKeysCommand` is run.

When `~/.ssh/authorized_keys` does not contain the keys required to
authenticate a user, `sshd` invokes `ussher` to provide additional,
remotely-sourced keys, such as from Github or another identity and access
management provider.

## How it works

`ussher` provides a fallback mechanism for statically-defined `authorized_keys`
files, such as when `authorized_keys` needs to be frequently updated, composed
from a large number of sources, or simply defined at the moment of
authorization.

1. When you `ssh $USER@$HOSTNAME`, `sshd` first reads something like
   `/home/.ssh/authorized_keys` to authenticate incoming SSH connections. If
   `authorized_keys` contains a public key that matches the incoming
   connection, then the connection attempt proceeds normally.

1. If that file does not exist or does not contain a public key for the
   incoming connection, `sshd` will invoke `AuthorizedKeysCommand` as
   `AuthorizedKeysCommandUser`. The `AuthorizedKeysCommand` (`ussher`, in this
   case), is responsible for returning a list of authorized public keys, which
   `sshd` then uses to validate the incoming connection.

1. When invoked, `ussher` sources authorized keys from any number of remote
   sources, such as a static text file or more commonly something like
   `https://github.com/{username}.keys`. `ussher` then returns a single set of
   authorized keys to `sshd` which are used to validate the incoming
   connection.

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
   AuthorizedKeysCommand /usr/bin/ussher
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

`/etc/ussher` contains configuration files for each user it supports. For
example, to allow @dolph to SSH to your host as root (but, you know, _don't_),
you would configure `/etc/ussher/root.yml` using:

```yaml
sources:
- url: https://github.com/dolph.keys
```
