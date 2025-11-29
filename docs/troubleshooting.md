# Troubleshooting

This guide covers common issues and their solutions.

## Installation Issues

### Homebrew Installation Fails

**Problem**: `brew install sshbuddy` fails with "formula not found"

**Solution**: Ensure you've tapped the repository first:
```bash
brew tap javedh-dev/tap
brew install sshbuddy
```

### Binary Won't Execute

**Problem**: Downloaded binary shows "permission denied" error

**Solution**: Make the binary executable:
```bash
chmod +x sshbuddy
```

On macOS, you may also need to allow the app in System Preferences > Security & Privacy.

## Configuration Issues

### Config File Not Found

**Problem**: SSHBuddy can't find or create the config file

**Solution**: Ensure the config directory exists and is writable:
```bash
mkdir -p ~/.config/sshbuddy
chmod 755 ~/.config/sshbuddy
```

### Invalid JSON Error

**Problem**: SSHBuddy reports "invalid JSON" on startup

**Solution**: Your config file may be corrupted. Back it up and validate the JSON:
```bash
cp ~/.config/sshbuddy/config.json ~/.config/sshbuddy/config.json.backup
cat ~/.config/sshbuddy/config.json | python -m json.tool
```

If validation fails, either fix the JSON manually or delete the file to let SSHBuddy create a fresh one.

## SSH Config Integration

### Hosts Not Appearing from SSH Config

**Problem**: Hosts defined in `~/.ssh/config` don't show up in SSHBuddy

**Solution**:
1. Press `s` to open settings
2. Verify "SSH Config" is enabled (green checkmark)
3. Check that your SSH config file exists: `ls -la ~/.ssh/config`
4. Ensure the config file is readable: `chmod 644 ~/.ssh/config`

### Custom SSH Config Path Not Working

**Problem**: Specified a custom SSH config path but hosts don't load

**Solution**:
1. Press `s`, navigate to "SSH Config", and press `e`
2. Verify the path is correct and absolute (e.g., `/home/user/.ssh/custom-config`)
3. Ensure the file exists and is readable
4. Leave the field empty to use the default `~/.ssh/config`

## Termix Integration

### Authentication Fails

**Problem**: "Authentication failed" error when connecting to Termix

**Solution**:
1. Verify your Termix server is accessible: `curl https://your-termix-server.com/api`
2. Check that the base URL is correct (should end with `/api`)
3. Ensure your username and password are correct
4. Check Termix server logs for authentication errors

### Token Expired Loop

**Problem**: Constantly prompted for credentials even after successful login

**Solution**: This usually indicates the Termix server isn't setting the JWT cookie properly. Check:
1. Termix server is returning a `Set-Cookie` header with `jwt` cookie
2. Cookie has a valid `MaxAge` or `Expires` value
3. Server time is synchronized (JWT expiry depends on accurate time)

### Hosts Not Loading from Termix

**Problem**: Authentication succeeds but no hosts appear

**Solution**:
1. Verify the `/ssh/db/host` endpoint returns data: 
   ```bash
   curl -H "Cookie: jwt=YOUR_TOKEN" https://your-termix-server.com/api/ssh/db/host
   ```
2. Check that the response is a JSON array of host objects
3. Ensure hosts have required fields: `name`, `ip`, `username`, `port`

## Display Issues

### Terminal Too Small Error

**Problem**: "Terminal Too Small" message appears

**Solution**: SSHBuddy requires a minimum terminal size of 84x24. Resize your terminal window or use a larger font size.

### Colors Look Wrong

**Problem**: Colors appear incorrect or hard to read

**Solution**:
1. Try a different theme (press `s`, navigate to Theme, press Space/Enter)
2. Ensure your terminal supports 256 colors
3. Check your terminal's color scheme settings

### Icons Not Displaying

**Problem**: Source icons (◆, ■, ▲) appear as boxes or question marks

**Solution**: Your terminal font may not support these Unicode characters. Most modern terminal fonts include them, but if yours doesn't:
1. Install a font with better Unicode support (e.g., Nerd Fonts, Fira Code)
2. The icons are purely visual—functionality isn't affected

## Connection Issues

### SSH Connection Fails

**Problem**: Selecting a host and pressing Enter doesn't connect

**Solution**:
1. Verify the host details are correct (hostname, user, port)
2. Test the connection manually: `ssh user@hostname`
3. Check that your SSH client is installed: `which ssh`
4. Review SSH key permissions if using identity files: `chmod 600 ~/.ssh/id_rsa`

### ProxyJump Not Working

**Problem**: Connections through bastion hosts fail

**Solution**:
1. Ensure your SSH client supports ProxyJump (OpenSSH 7.3+)
2. Verify the bastion host is accessible
3. Test the ProxyJump manually: `ssh -J bastion.example.com user@target.example.com`

## Performance Issues

### Slow Startup

**Problem**: SSHBuddy takes a long time to launch

**Solution**:
1. Check if Termix integration is enabled and the server is slow/unreachable
2. Temporarily disable Termix to isolate the issue
3. Large SSH config files can slow parsing—consider splitting them

### Ping Takes Too Long

**Problem**: Pressing `p` to ping hosts is very slow

**Solution**: Pinging is done sequentially. If you have many hosts or some are unreachable with long timeouts, it will take time. This is normal behavior—pinging continues in the background while you use the interface.

## Debug Logs

SSHBuddy writes debug information to `/tmp/sshbuddy-debug.log`. Check this file for detailed error messages:

```bash
tail -f /tmp/sshbuddy-debug.log
```

This is especially useful for diagnosing Termix integration issues.

## Getting Help

If you're still experiencing issues:

1. Check the [GitHub Issues](https://github.com/javedh-dev/sshbuddy/issues) for similar problems
2. Review the debug log for error messages
3. Open a new issue with:
   - Your operating system and version
   - SSHBuddy version (`sshbuddy --version`)
   - Steps to reproduce the problem
   - Relevant log entries (sanitize any sensitive information)

## Reporting Bugs

When reporting bugs, please include:
- Operating system and version
- Terminal emulator and version
- SSHBuddy version
- Steps to reproduce
- Expected vs actual behavior
- Debug log excerpts (remove sensitive data)

The more information you provide, the easier it is to diagnose and fix the issue.
