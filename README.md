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

## Configuration

`opt/ussher` contains configuration files for each user it supports. For example, to allow me to SSH to your host as root, you would use:

`/opt/ussher/root.yml`:

```yaml
urls:
- https://github.com/dolph.keys
```

## Recommended installation & usage

1. Create a dedicated user and group to run `ussher`, named `ussher`:

   ```bash
   sudo useradd --system --user-group ussher
   ```

2. Install the `ussher` binary to `/usr/sbin/ussher`.

   ```bash
   sudo cp ussher /usr/bin/ussher
   sudo chown root:root /usr/bin/ussher
   sudo chmod 0700 /usr/bin/ussher
   ```

3. Create a directory for caching remotely-sourced data.

   ```bash
   sudo mkdir --parents /var/cache
   sudo mkdir --parents --mode=0700 /var/cache/ussher
   sudo chown ussher:ussher /var/cache/ussher
   ```

4. Create a directory for logging.

   ```bash
   sudo mkdir --parents --mode=0700 /var/log/ussher
   sudo chown ussher:ussher /var/log/ussher
   ```

5. Configure `sshd` to invoke `ussher`. Add the following lines to
   `/etc/ssh/sshd_config`:

   ```
   AuthorizedKeysCommand /usr/bin/ussher
   AuthorizedKeysCommandUser ussher
   ```

   Restart `sshd`, for example:

   ```bash
   sudo systemctl restart sshd
   ```
