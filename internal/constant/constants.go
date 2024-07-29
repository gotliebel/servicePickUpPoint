package constant

const (
	DataBaseConnection   = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	OrderStorageFilename = "orders.json"
	TimeFormat           = "02-01-2006"
	Success              = "Success!"
	WrongCmdTextTemplate = "Command %s doesn't exist. Try command help\n"
	WaitForCmdTxt        = "Enter your command:"
)

const TakebackPageSize = 10

const (
	FlagOrderId     = "order_id"
	FlagClientId    = "client_id"
	FlagStoredUntil = "stored_until"
	FlagNumRoutines = "num_routines"
	FlagPack        = "package"
	FlagWeight      = "weight"
	FlagPrice       = "price"
)

const (
	CommandHelp              = "help"
	CommandAccept            = "accept"
	CommandBack              = "back"
	CommandPickUp            = "pickup"
	CommandList              = "list"
	CommandReturn            = "return"
	CommandTakebacks         = "takebacks"
	CommandExit              = "exit"
	CommandChangeNumRoutines = "change"
)

const (
	CustomExit = "routine_exit"
)

const HelpTxt = `This is a command line tool for pickup point
Currently these commands available:
	` +
	CommandHelp + `: which prints this entire text
	` +
	CommandExit + `: ends work of programm
	` +
	CommandAccept + `: to accept an order from courier
		flags:
		-` + FlagOrderId + ` (positive integer)
		-` + FlagClientId + ` (positive integer)
		-` + FlagStoredUntil + ` (date in dd-mm-yyyy format)
		-` + FlagPack + ` (type of package)
		-` + FlagWeight + ` (weight of order)
		-` + FlagPrice + ` (price of order)
	` +
	CommandBack + `: to give back to courier orders, whose storage time is over
		flags:
		-` + FlagClientId + ` (positive integer)
	` +
	CommandPickUp + `: to pickup client's order
		flags:
        	-list (a string of id's delimited by ',', for example -list=1,2,3,144)
	` +
	CommandList + `: to give a list of client's orders
		flags:
		-` + FlagClientId + ` (positive integer)
        	[-limit (positive integer)]
	` +
	CommandReturn + `: to accept a return from client
		flags:
		-` + FlagOrderId + ` (positive integer)
		-` + FlagClientId + ` (positive integer)
	` +
	CommandTakebacks + `: to get a page from list of takebacks. Page size is 10
		flags:
			-page(positive integer)
		`
