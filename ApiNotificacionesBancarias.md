# Api Notificaciones Bancarias

# Endpoints

- Generalidades
    
    Todas las peticiones a los módulos de los bancos requieren una API KEY independientemente de la naturaleza de la petición.
    
    Para cada Banco solo puede existir una API Key.
    
    ```json
    {
    	"X-API-Key" :"757d488b9ab8aeb70388bd0213b877be2d4c5918873afab4b7227fd3add1248b"
    }
    
    ```
    
- “/”
    
    Endpoint para verificar si el servidor esta activo.
    
    Se espera un código:  **200**
    
    ```json
    {
    	"message": "api notificaciones UP!"
    }
    ```
    

## BDV

- **Weebhook BDV**
    
    **Endpoint:** {URL}/api-notificaciones/v1/bdv/weebhook
    
    **Tipo de petición:** POST
    
    Enpoint WeebHook para Banco de Venezuela, envía la siguiente estructura de datos:
    
    ```json
    	{
    		"bancoOrdenante": "0191",
    		"referenciaBancoOrdenante": "555667",
    		"numeroCliente": "123456",
    		"numeroClientemercio": "12346",
    		"fecha": "2023-10-10",
    		"hora": "23:11",
    		"monto": "100",
    		"idCliente": "V20697579",
    		"numeroComercio": "04244648106",
    		"idComercion": "V20697579"
    	}
    ```
    
    Lo campos son todos strings.
    
    BDV enviara constantemente los pagos recibidos a la cuenta de la empresa a este WeebHook.
    
    ## Respuesta 201:
    
    Un código 201 implica que una entrada fue creada exitosamente y es el codigo esperado por BDV cuando se crea la entrada de un pago.
    
    El **codigo** es un campo de confirmación que espera BDV. **mensajeSistema** y **mensajeCliente** igualmente son campos esperados por BDV.
    
    **success** es un campo interno que sirve como flag para evaluar rápidamente la situación de la respuesta.
    
    El siguiente ejemplo muestra los valores que esperan para un 201:
    
    ```json
    {
    			"codigo":         "00",
    			"mensajeCliente": "Aprobado",
    			"mensajeSistema": "Notificado",
    			"success":        result,
    }
    ```
    
    ## Respuestas 200:
    
    Un código 200 de utiliza cuando un pago ya fue registrado con anterioridad, espera los mismos campos pero difiere en valores.
    
    ```json
    {
    		"codigo":         "01",
    		"mensajeCliente": "pago previamente recibido",
    		"mensajeSistema": "Notificado",
    		"success": true
    }
    ```
    
    ## Respuestas 400:
    
    Las respuestas 400 implican incumplimiento de la convención esperada en los campos de la petición.
    
    En estos caso el **codigo** es null y success es false.
    
    Los campos de mensaje indican el proceso que fallo y detalles del error, al ocurrir se puede ir al código fuente para verificar la razon del error de manera mas detallada.
    
    ```json
    {
    	"codigo": null,
    	"mensajeCliente": "validation error",
    	"mensajeSistema": "date es obligatorio",
    	"success": false
    }
    ```
    
- **Obtención de pagos BDV**
    
    **Endpoint**: {URL}/api-notificaciones/v1/api-notificaciones/v1/bdv/notificaciones?referencia={value}&fecha={value}&numero_cliente{value}&id_cliente={value}&pagina={value}
    
    Este endpoint permite recuperar y verificar los pagos recibidos. 
    Los valores del las querys son opcionales.
    
    Si no si se proporcionan valores retronara los primeros 100 valores en la primera pagina.
    
    Requiere uso de API Key, solo puede acceder el BDV y Desarrollo.
    
    ```json
    {
    	"data": [
    		{
    			"BancoOrigen": "0191",
    			"ReferenciaOrigen": "555667",
    			"IdCliente": "V20697579",
    			"NumeroCliente": "123456",
    			"IdComercio": "04244648106",
    			"NumeroComercio": "",
    			"FechaBanco": "2023-10-10",
    			"FechaTranformada": "2023-10-10T00:00:00Z",
    			"HoraBanco": "23:11",
    			"HoraTransformada": "0000-01-01T23:11:00Z",
    			"Monto": 100
    		}
    	],
    	"message": "Notifications retrieved successfully",
    	"pagination": {
    		"current_page": 1,
    		"page_size": 100,
    		"total_items": 1,
    		"total_pages": 1
    	},
    	"success": true
    }
    ```
    
    El campo **data** siempre es una lista/arreglo/slice de los datos obtenidos, se corresponden con el modelo **NotificationBDV**, **current_page** la pagina actual de los datos, **page_size** el tamaño de la pagina, **total_items** la cantidad de items encontrado en la pagina actual y **total_pages** la cantidad de paginas calculadas según el tamaño de las paginas.
    

## Bancaribe

- **Weebhook Bancaribe**
    
    **Endpoint:** {URL}/api-notificaciones/v1/bancaribe/weebhook
    
    **Tipo de petición**: POST
    
    **Descripción:** Endpoint WeebHook para Bancaribe, envía la siguiente estructura de datos:
    
    ```json
    {
    	 "amount": 100,
    	 "bankName": "BANCO DEL CARIBE",
    	 "clientPhone": "00584247776589",
    	 "commercePhone": "00584168327199",
    	 "creditorAccount": "01140152001520123861",
    	 "currencyCode": "VES",
    	 "date": "23-10-2024",
    	 "debtorAccount": "01140152001520123746",
    	 "debtorID": "411823643",
    	 "destinyBankReference": "000254151380",
    	 "originBankCode": "0114",
    	 "originBankReference": "254151380",
    	 "paymentType": "TRF",
    	 "time": "08:45:00"
    }
    ```
    
    A diferencia de BDV, se espera que los montos lleguen como números, ellos proporcionan un numero entero pero asumimos que pueden ser flotantes.
    
    Bancaribe enviara constantemente los pagos recibidos a la cuenta de la empresa a este WeebHook.
    
    ## Respuesta 201:
    
    Un código 201 implica que una entrada fue creada exitosamente y es el codigo esperado por BDV cuando se crea la entrada de un pago.
    
    El **codigo** es un campo de confirmación que espera BDV. **mensajeSistema** y **mensajeCliente** igualmente son campos esperados por BDV.
    
    **success** es un campo interno que sirve como flag para evaluar rápidamente la situación de la respuesta.
    
    El siguiente ejemplo muestra los valores que esperan para un 201:
    
    ```json
    {
    	"codigo": "00",
    	"data": {
    		"ID": 6,
    		"Amount": 100,
    		"BankName": "BANCO DEL CARIBE",
    		"ClientPhone": "00584247776589",
    		"CommercePhone": "00584168327199",
    		"CreditorAccount": "01140152001520123861",
    		"CurrencyCode": "VES",
    		"DateBancaribe": "23-10-2024",
    		"Date": "2024-10-23T00:00:00Z",
    		"DebtorID": "411823643",
    		"DestinyBankReference": "000254851388",
    		"OriginBankCode": "0114",
    		"OriginBankReference": "254151380",
    		"PaymentType": "TRF",
    		"TimeBancaribe": "08:45:00",
    		"Time": "0000-01-01T08:45:00Z"
    	},
    	"message": "Success",
    	"messageSystem,": "Notificado",
    	"statusCode": 200,
    	"success": true
    }
    ```
    
    ## Respuestas 200:
    
    Un código 200 de utiliza cuando un pago ya fue registrado con anterioridad, espera los mismos campos pero difiere en valores.
    
    ```json
    {
    	"codigo": "01",
    	"message": "pago previamente recibido",
    	"messageSystem,": "Notificado",
    	"statusCode": 200,
    	"success": true
    }
    ```
    
    ***Ver nota de** 
    
    ## Respuestas 400:
    
    Las respuestas 400 implican incumplimiento de la convención esperada en los campos de la petición.
    
    En estos caso el **codigo** es null y success es false.
    
    Los campos de mensaje indican el proceso que fallo y detalles del error, al ocurrir se puede ir al código fuente para verificar la razon del error de manera mas detallada.
    
    ```json
    {
    	"codigo": null,
    	"message": "validation error",
    	"messageSystem,": "paymentType es obligatorio",
    	"statusCode": 400,
    	"success": false
    }
    ```
    
- **Obtención de pagos Bancaribe**
    
    **Endpoint:** {URL}/api-notificaciones/v1/api-notificaciones/v1/bancaribe/notificaciones?referencia={value}&fecha={value}&numero_cliente{value}&id_cliente={value}&pagina={value}
    
    Este endpoint permite recuperar y verificar los pagos recibidos.
    Los valores del las querys son opcionales.
    
    Si no si se proporcionan valores retronara los primeros 100 valores en la primera pagina.
    
    Requiere uso de API Key, solo puede acceder el Bancaribe y Desarrollo.
    
    ```json
    {
    	"data": [
    		{
    			"ID": 1,
    			"Amount": 100,
    			"BankName": "BANCO DEL CARIBE",
    			"ClientPhone": "00584247776589",
    			"CommercePhone": "00584168327199",
    			"CreditorAccount": "01140152001520123861",
    			"CurrencyCode": "VES",
    			"DateBancaribe": "23-10-2024",
    			"Date": "2024-10-23T00:00:00Z",
    			"DebtorID": "411823643",
    			"DestinyBankReference": "000254851388",
    			"OriginBankCode": "0114",
    			"OriginBankReference": "254151380",
    			"PaymentType": "TRF",
    			"TimeBancaribe": "08:45:00",
    			"Time": "0000-01-01T08:45:00Z"
    		}
    	],
    	"message": "Notifications retrieved successfully",
    	"pagination": {
    		"current_page": 1,
    		"page_size": 100,
    		"total_items": 1,
    		"total_pages": 1
    	},
    	"success": true
    }
    ```
    

## Tesoro

En construcción

# Models

## Generalidades

Los campos en la BD se tratan de mantener con los mismos nombres con los que son recibidos, lo que para algunos modelos estarán en ingles y para otros en español.

Tanto en las columnas de la BD como en el JSON de retorno de la información se mantienen en snake_case, siendo este la convención de Go (mas no de obligatoriedad).

Ha diferencia de otros sistema, aquí los ID si los mantenemos como números enteros dado que no esperamos tener distintas instancias de la misma API, además Go y Postgres deberían ser suficientemente rápidos al momento de recibir los datos y realizar las escrituras pertinentes.

## BDV

**Modelo:**

```go
type NotificationBDV struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	BancoOrigen      string    `gorm:"column:banco_origen;size:4" json:"banco_origen"`                // banco del cliente
	ReferenciaOrigen string    `gorm:"column:referencia_origen;size:15" json:"referencia_origen"`     // referencia del banco del cliente
	IdCliente        string    `gorm:"column:id_cliente;size:15" json:"id_cliente"`                   // CI/Rif cliente
	NumeroCliente    string    `gorm:"column:numero_cliente;size:15" json:"numero_cliente"`           // tlf cliente
	IdComercio       string    `gorm:"column:id_comercio;size:23" json:"id_comercio"`                 // Rif Comercio
	NumeroComercio   string    `gorm:"column:numero_comercio;size:15" json:"numero_comercio"`         // Tlf Comercio
	FechaBanco       string    `gorm:"column:fecha_banco;size:11" json:"fecha_banco"`                 // Fecha en str
	FechaTranformada time.Time `gorm:"column:fecha_transformada;type:date" json:"fecha_transformada"` // transformar para crear busquedas
	HoraBanco        string    `gorm:"column:hora_banco;size:7" json:"hora_banco"`                    // hora en str
	HoraTransformada time.Time `gorm:"column:hora_transformada;type:time" json:"hora_transformada"`
	Monto            float64   `gorm:"column:monto;type:decimal(13,2)" json:"monto"` // previendo futuras conversiones monetarias
}
```

**Campos:**

- **BancoOrigen** deben ser siempre cuatro dígitos.
- **ReferenciaOrigen:** numero de referencia emitido por el banco que la origina.
- **IdCliente:** CI/RIF del cliente pagador.
- **NumeroCliente:** numero de teléfono del cliente pagador.
- **IdComercio:** RIF del cliente receptor (Netcom Plus).
- **NumeroComercio:** numero de teléfono del cliente receptor (Netcom Plus).
- **FechaBanco:** fecha en string tal cual llega desde el banco.
- **FechaTransformada:** fecha llevado a un objeto Time.
- **HoraBanco:** hora en string tal cual la entrega el Banco.
- **HoraTransformada:** hora llevado a objeto time tal cual la entrega el banco.
- **Monto:** valor del pago como punto flotante.

## Bancaribe

**Modelo:**

```go
type NotificationBancaribe struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Amount               float64   `gorm:"type:decimal(13,2);column:amount" json:"amount"`
	BankName             string    `gorm:"size:30;column:bank_name" json:"bank_name"`
	ClientPhone          string    `gorm:"size:16;column:client_phone" json:"client_phone"`
	CommercePhone        string    `gorm:"size:16;column:commerce_phone" json:"commerce_phone"`
	CreditorAccount      string    `gorm:"size:50;column:creditor_account" json:"creditor_account"`
	CurrencyCode         string    `gorm:"size:5;column:currency_code" json:"currency_code"`
	DateBancaribe        string    `gorm:"size:12;column:date_bancaribe" json:"date_bancaribe"`
	Date                 time.Time `gorm:"type:date;column:date" json:"date"`
	DebtorID             string    `gorm:"size:15;column:debtor_id" json:"debtor_id"`
	DestinyBankReference string    `gorm:"size:15;column:destiny_bank_reference" json:"destiny_bank_reference"`
	OriginBankCode       string    `gorm:"size:5;column:origin_bank_code" json:"origin_bank_code"`
	OriginBankReference  string    `gorm:"size:15;column:origin_bank_reference" json:"origin_bank_reference"`
	PaymentType          string    `gorm:"size:6;column:payment_type" json:"payment_type"`
	TimeBancaribe        string    `gorm:"size:10;column:time_bancaribe" json:"time_bancaribe"`
	Time                 time.Time `gorm:"type:time;column:time" json:"time"`
}

```

**Campos:**

- **Amount:** monto del pago recibido, se procesa como punto flotante de hasta dos dígitos.
- **BankName:** nombre del Banco Pagador.
- **ClientPhone:** teléfono del cliente pagador.
- **CommercePhone:** teléfono del cliente receptor.
- **CreditorAccount:** Cuenta del cliente pagador
- **CurrencyCode:** código de moneda de pago.
- **DateBancaribe:** fecha del pago en formato string.
- **Date:** fecha en formato fecha.
- **DebtorID:** CI/RIF del cliente pagador.
- **DestiniBankReference:** referencia otorgada al pago por Bancaribe.
- **OriginBankCode:** código numérico de 4 dígitos del banco.
- **OriginBankReference:** referencia del banco pagador.
- **PaymentType:** tipo de pago (TRF: transferencia, PM: Pagomóvil)
- **TimeBancaribe:** hora del pago en formato string.
- **Time:** hora del pago en formato Time.

## Tesoro

En construcción  🏗️

# Handlers