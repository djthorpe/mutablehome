/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

// https://github.com/ggravlingen/pytradfri/blob/master/pytradfri/const.py
const (
	ROOT_DEVICES                              = "15001"
	ROOT_GATEWAY                              = "15011"
	ROOT_GROUPS                               = "15004"
	ROOT_MOODS                                = "15005"
	ROOT_NOTIFICATION                         = "15006" // speculative name
	ROOT_REMOTE_CONTROL                       = "15009"
	ROOT_SIGNAL_REPEATER                      = "15014"
	ROOT_SMART_TASKS                          = "15010"
	ROOT_START_ACTION                         = "15013" // found under ATTR_START_ACTION
	ATTR_START_BLINDS                         = "15015"
	ATTR_ALEXA_PAIR_STATUS                    = "9093"
	ATTR_AUTH                                 = "9063"
	ATTR_APPLICATION_TYPE                     = "5750"
	ATTR_BLIND_CURRENT_POSITION               = "5536"
	ATTR_BLIND_TRIGGER                        = "5523"
	ATTR_CERTIFICATE_PEM                      = "9096"
	ATTR_CERTIFICATE_PROV                     = "9092"
	ATTR_CLIENT_IDENTITY_PROPOSED             = "9090"
	ATTR_CREATED_AT                           = "9002"
	ATTR_COGNITO_ID                           = "9101"
	ATTR_COMMISSIONING_MODE                   = "9061"
	ATTR_CURRENT_TIME_UNIX                    = "9059"
	ATTR_CURRENT_TIME_ISO8601                 = "9060"
	ATTR_DEVICE_INFO                          = "3"
	ATTR_GATEWAY_ID_2                         = "9100" // stored in IKEA app code as gateway id
	ATTR_GATEWAY_TIME_SOURCE                  = "9071"
	ATTR_GATEWAY_UPDATE_PROGRESS              = "9055"
	ATTR_GROUP_MEMBERS                        = "9018"
	ATTR_HOMEKIT_ID                           = "9083"
	ATTR_HS_LINK                              = "15002"
	ATTR_ID                                   = "9003"
	ATTR_IDENTITY                             = "9090"
	ATTR_IOT_ENDPOINT                         = "9103"
	ATTR_KEY_PAIR                             = "9097"
	ATTR_LAST_SEEN                            = "9020"
	ATTR_LIGHT_CONTROL                        = "3311" // array
	ATTR_MASTER_TOKEN_TAG                     = "9036"
	ATTR_MOOD                                 = "9039"
	ATTR_NAME                                 = "9001"
	ATTR_NTP                                  = "9023"
	ATTR_FIRMWARE_VERSION                     = "9029"
	ATTR_FIRST_SETUP                          = "9069" // ??? unix epoch value when gateway first setup
	ATTR_GATEWAY_INFO                         = "15012"
	ATTR_GATEWAY_ID                           = "9081" // ??? id of the gateway
	ATTR_GATEWAY_REBOOT                       = "9030" // gw reboot
	ATTR_GATEWAY_FACTORY_DEFAULTS             = "9031" // gw to factory defaults
	ATTR_GATEWAY_FACTORY_DEFAULTS_MIN_MAX_MSR = "5605"
	ATTR_GOOGLE_HOME_PAIR_STATUS              = "9105"
	ATTR_DEVICE_STATE                         = "5850" // 0 / 1
	ATTR_LIGHT_DIMMER                         = "5851" // Dimmer, not following spec: 0..255
	ATTR_LIGHT_COLOR_HEX                      = "5706" // string representing a value in hex
	ATTR_LIGHT_COLOR_X                        = "5709"
	ATTR_LIGHT_COLOR_Y                        = "5710"
	ATTR_LIGHT_COLOR_HUE                      = "5707"
	ATTR_LIGHT_COLOR_SATURATION               = "5708"
	ATTR_LIGHT_MIREDS                         = "5711"
	ATTR_NOTIFICATION_EVENT                   = "9015"
	ATTR_NOTIFICATION_NVPAIR                  = "9017"
	ATTR_NOTIFICATION_STATE                   = "9014"
	ATTR_OTA_TYPE                             = "9066"
	ATTR_OTA_UPDATE_STATE                     = "9054"
	ATTR_OTA_UPDATE                           = "9037"
	ATTR_PUBLIC_KEY                           = "9098"
	ATTR_PRIVATE_KEY                          = "9099"
	ATTR_PSK                                  = "9091"
	ATTR_REACHABLE_STATE                      = "9019"
	ATTR_REPEAT_DAYS                          = "9041"
	ATTR_SEND_CERT_TO_GATEWAY                 = "9094"
	ATTR_SEND_COGNITO_ID_TO_GATEWAY           = "9095"
	ATTR_SEND_GH_COGNITO_ID_TO_GATEWAY        = "9104"
	ATTR_SENSOR                               = "3300"
	ATTR_SENSOR_MAX_RANGE_VALUE               = "5604"
	ATTR_SENSOR_MAX_MEASURED_VALUE            = "5602"
	ATTR_SENSOR_MIN_RANGE_VALUE               = "5603"
	ATTR_SENSOR_MIN_MEASURED_VALUE            = "5601"
	ATTR_SENSOR_TYPE                          = "5751"
	ATTR_SENSOR_UNIT                          = "5701"
	ATTR_SENSOR_VALUE                         = "5700"
	ATTR_START_ACTION                         = "9042" // array
	ATTR_SMART_TASK_TYPE                      = "9040" // 4 = transition | 1 = not home | 2 = on/off
	ATTR_SMART_TASK_NOT_AT_HOME               = 1
	ATTR_SMART_TASK_LIGHTS_OFF                = 2
	ATTR_SMART_TASK_WAKE_UP                   = 4
	ATTR_SMART_TASK_TRIGGER_TIME_INTERVAL     = "9044"
	ATTR_SMART_TASK_TRIGGER_TIME_START_HOUR   = "9046"
	ATTR_SMART_TASK_TRIGGER_TIME_START_MIN    = "9047"
	ATTR_SWITCH_CUM_ACTIVE_POWER              = "5805"
	ATTR_SWITCH_ON_TIME                       = "5852"
	ATTR_SWITCH_PLUG                          = "3312"
	ATTR_SWITCH_POWER_FACTOR                  = "5820"
	ATTR_TIME_END_TIME_HOUR                   = "9048"
	ATTR_TIME_END_TIME_MINUTE                 = "9049"
	ATTR_TIME_START_TIME_HOUR                 = "9046"
	ATTR_TIME_START_TIME_MINUTE               = "9047"
	ATTR_TRANSITION_TIME                      = "5712"
	ATTR_USE_CURRENT_LIGHT_SETTINGS           = "9070"
)
