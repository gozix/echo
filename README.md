# GoZix Echo

## Dependencies

* [validator](https://github.com/gozix/validator)
* [viper](https://github.com/gozix/viper)
* [zap](https://github.com/gozix/zap)

## Configuration

```json
{
  "echo": {
    "debug": false,
    "level": "debug",
    "static": {
      "prefix": "/",
      "root": ""
    },
    "hide_banner": false,
    "hide_port": false 
  }
}
```

## Built-in Tags

| Symbol                | Value                         | Description               | 
| --------------------- | ----------------------------- | ------------------------- |
| TagController         | echo.controller               | Add a controller          |
| TagConfigurator       | echo.configurator             | Add a configurator        |
| TagMiddleware         | echo.middleware               | Add a middleware          |
