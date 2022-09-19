# digitalSTROM

digitalSTROM (dS) connects electrical devices via the existing power line and enables the implementation of control and regulation tasks. This library connects either to the local ds-Server JSON-API or via Kitepage link in order to control devices or receive their sensor data. 

A complete documentation of the local API could be found [here](https://developer.digitalstrom.org/Architecture/v1.7/).

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

# Understanding digitalSTROM local API (dSS)

## Account

The module ```Account``` is clustering everything to communicate with one single local JSON API. It caches and updates Devices, Sensor values, Channel values, etc for an easier access. 

## Circuits / Metering

Circuits or Metering are the representation of the digitalSTROM meters. They measure the consumption and the used energy for multiple devices. Each devices is assigned to one meter. 
However, not all Circuits are able to do metering. ```HasMetering``` is a flag to validate whether metering data could be requested via the API. 

## Structure
The Structure object covers all components related to the digitalSTROM communication network in the smart home. The logical structure in the smart home consists of Apartment and Zones. Zones contain Devices and Applications (Groups). Devices could be assigned to multiple groups.


     Structure ┐
               └ Apartment ┐
                           ├ Zones ┐
                           |       ├ Devices ┐
                           |       |         ├ Sensors
                           |       |         ├ BinaryInputs
                           |       |         ├ OutputChannels
                           |       |         └ GroupIds
                           |       └ Groups ┐
                           |                └ DeviceIDs
                           └ Floors ┐
                                    └ ZoneIDs

## Device

     Device ┐
            ├ Sensors
            ├ BinaryInputs
            ├ OutputChannels
            └ GroupIds

### Sensors


### Binary Inputs

### OutputChannels

Channel values could be read with the function 

    Account.PollChannelValue(channel *OutputChannel) (int, error)

and set with the function

    Account.SetOutputChannelValue(channel *OutputChannel, value string)

Keep in mind: When reading a channel value like the brightness of a lamp, the value range is 0-255 whereas 255 means 100%. However, when setting the value, the range is 0-100. The library is not compansating this issue and simply forwards the values. Example, if you setg the brightness to 50% (value 50), the read value will then be 128.

### ON State

On states can't be requested directly via the API. Instead a complete structure request has to be performed. The function 

    Account.PollStructureValues()

will perform a structure requests and updates all ```On``` states for all devices as well as there ```IsPresent``` information. Depending on the device type, the ```On``` state might have different behaviors (see "ON State for Lamps")

### ON State for Lamps

Setting the on state by the command 

    Account.TurnOn(device *Device, on bool) 

the lamp will turn or off (depending on the parameter on). However, by setting the output channel value for ```brightness``` to a value higher than 0% and does not force the system to set the ```On``` state to be set to ```true```. Same vise versa, setting the ```brightness``` to 0% will not set the ```On``` value to ```false```. 
