# digitalSTROM

digitalSTROM (dS) connects electrical devices via the existing power line and enables the implementation of control and regulation tasks. This library connects either to the local ds-Server JSON-API or via Kitepage link in order to control devices or receive their sensor data. 

# How to use

## Local access

### Create new account instance

    account := *digitalstrom.NewAccount()

### Set URL

    account.SetURL("local or pagekite link including protocol and port")

### Register an Application

To avoid a handling with ``userName`` and ``password``, each app could register at the dS server in order to receive an application token. This token will be used to perform an application login. This library requires the application token in order to work. Register an application only once. The application token stays valid until the user deletes the access rights manually at the dS server side. 

    token, err := account.RegisterApplication("username","userpassw","application name")

### Application Login


    account.SetApplicationToken("your token")


### Inititialization

The initialization processes several tasks. An application login will be performed in order
to receive the ``session token``. This token will kept in memory and used for future requests.
After a successful login, the complete structure and circuits will be requested. The ``session token```
will be refreshed automatically.


    err := account.Init()
