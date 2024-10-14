
# wait-for-it

`wait-for-it` is a script to wait for multiple TCP host:port combinations to become available, useful for orchestrating dependent services in distributed systems or containerized environments.

## Features

- **Supports Multiple Hosts**: Monitor multiple host:port pairs at the same time.
- **Timeout and Retry**: Allows configuring timeout and retry intervals.
- **Flexible Log Formatting**: Logs can be output in either plain text or JSON format.
- **Environment Variable Support**: Optionally pass host:port pairs via the `WAIT` environment variable.
- **Log Verbosity**: Control log verbosity via flags (quiet, debug).

## Usage

You can either pass the host:port pairs as positional arguments or by using the `WAIT` environment variable.


### Example Command:

```bash
./wait-for-it -timeout 30 host1:8080,host2:9090
```

### Using Environment Variables:

You can also set the `WAIT` environment variable to provide the host:port pairs:

```bash
WAIT=host1:8080,host2:9090 ./wait-for-it -timeout 30
```

If both the environment variable and positional arguments are provided, the positional arguments take precedence.

## Options

| Option               | Description                                                                                                                                                                  |
|----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-timeout`           | Set timeout in seconds (default: 15).                                                                                                                                        |
| `-retry-interval`     | Set retry interval in seconds (default: 1).                                                                                                                                  |
| `-quiet`             | Enable quiet mode (suppress output, only warnings and errors).                                                                                                               |
| `-debug`             | Enable debug mode (more detailed output).                                                                                                                                    |
| `-format`            | Set log format (options: 'text' or 'json') (default: 'text').                                                                                                                |
| `-help`              | Show help information.                                                                                                                                                       |

## Environment Variables

| Variable | Description                                                                                                                                                      |
|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `WAIT`   | A comma-separated list of host:port pairs, e.g. `host1:8080,host2:9090`. This can be used instead of passing host:port pairs as positional arguments.             |

## Example Usage

### 1. Check if multiple services are ready with a 30-second timeout:
```bash
./wait-for-it -timeout 30 host1:8080,host2:9090
```

### 2. Use the `WAIT` environment variable for host:port pairs:
```bash
WAIT=host1:8080,host2:9090 ./wait-for-it -timeout 30
```

### 3. Check with a custom retry interval of 5 seconds between checks:
```bash
./wait-for-it -timeout 30 -retry-interval 5 host1:8080,host2:9090
```

### 4. Enable debug mode to see detailed logs for each connection attempt:
```bash
./wait-for-it -timeout 30 -debug host1:8080,host2:9090
```

### 5. Use quiet mode to suppress all output except warnings and errors:
```bash
./wait-for-it -quiet -timeout 30 host1:8080,host2:9090
```

### 6. Output logs in JSON format:
```bash
./wait-for-it -timeout 30 -format json host1:8080,host2:9090
```

## License

This project is open-source and available under the MIT License.
