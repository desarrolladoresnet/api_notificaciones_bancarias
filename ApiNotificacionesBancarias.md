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
    

---

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
    

---

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
    

---

## Tesoro

En construcción

---

# Models

## Generalidades

Los campos en la BD se tratan de mantener con los mismos nombres con los que son recibidos, lo que para algunos modelos estarán en ingles y para otros en español.

Tanto en las columnas de la BD como en el JSON de retorno de la información se mantienen en snake_case, siendo este la convención de Go (mas no de obligatoriedad).

Ha diferencia de otros sistema, aquí los ID si los mantenemos como números enteros dado que no esperamos tener distintas instancias de la misma API, además Go y Postgres deberían ser suficientemente rápidos al momento de recibir los datos y realizar las escrituras pertinentes.

---

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

---

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

---

## Tesoro

En construcción  🏗️

---

# Handlers

Son la funciones que manejan las peticiones recibidas a la API.

### BDV

### weebhook.go 🗃️

- **Función:** WeebHookBDV (Función Principal)
    
    **PROPOSITO**:
    
    Recibir los datos enviados desde BDV y retornar la respuesta esperada.
    
    **PROCESO:**
    
    1. Se recibe la notificación ⇒ Biding a un struct **bdvRequest:**
        1. Si ok, continua
        2. Not ok, retorna un 400 con success FALSE
    2. Valida que los campos no lleguen vacíos usando la función **Validate**.
        1. Si ok, continua
        2. Not ok, retorna un 400 con success FALSE
    3. Se transforma el **bdvRequest** en un  modelo **NotificationBDV**
        1. Si ok, continua
        2. Not ok, retorna un 400 con success FALSE
    4. Se verifica si la notificación existe previamente:
        1. Si existe, retorna un 200 con código “01” (exigencia de BDV).
        2. No existe, continua el proceso.
        3. Si hay un error al acceder a la BD, retorna un 500.
    5. Se guarda la notificación en la base de datos:
        1. Si ok, retorna un 201 con código “00” (exigencia de BDV).
        2. Error al escribir en la base de datos, retorna un 500.
        3. 
- **Función:** tranformRequestToModel
    
    **PROPOSITO:** tomar el struct de la petición y transformarlo en un struct que corresponda al modelo ***NotificationBDV.***
    
    **RETORNA:**
    
    1. NotificationBDV
    2. error
    
    **PROCESO:**
    
    1. Recepción del **bdvRequest**
    2. Transforma la fecha recibida en string en un objeto time.Time mediante la funci**ó**n **TransformDate.**
        1. Si todo ok, continua el proceso.
        2. Si falla retorna la notificación como **nil** y el error arrojado por la función
    3. Transforma la hora recibida en string en un objeto time.Time mediante la funci**ó**n **TransformHour.**
        1. Si todo ok, continua el proceso.
        2. Si falla retorna la notificación como **nil** y el error arrojado por la función
    4. Parsing del monto de string a float usando la ParseFloat  de la librería estándar de strconv.
        1. Si todo ok, continua el proceso.
        2. Si falla retorna la notificación como **nil** y el error arrojado por la función
    5. Se evalúa que el monto no sea menor a cero.
        1. Si todo ok, continua el proceso.
        2. Si falla retorna la notificación como **nil** y el error arrojado por la función
    6. Se crea el modelo y se retorna el puntero del mismo, error se retorna como **nil.**
- **Función:** TransformDate
    
    **PROPOSITO:** recibe el string con la fecha y la transforma en un objeto time***.***
    
    **RETORNA:**
    
    1. Time
    2. error
    
    **PROCESO:**
    
    ```go
    func TransformDate(date string) (*time.Time, error) {
    	parseDate, err := time.Parse("2006-01-02", date)
    	if err != nil {
    		return nil, err
    	}
    
    	return &parseDate, nil
    }
    ```
    
     
    
- **Función:** TransformHour
    
    **PROPOSITO:** recibe el string con la hora y la transforma en un objeto time***.***
    
    **RETORNA:**
    
    1. Time
    2. error
    
    **PROCESO:**
    
    ```go
    func TransformHour(timeStr string) (*time.Time, error) {
    	// Intentar con diferentes formatos
    	layouts := []string{"15.04", "15:04", "1504", "15 04"}
    
    	for _, layout := range layouts {
    		t, err := time.Parse(layout, timeStr)
    		if err == nil {
    			// Verificar rangos si el parseo fue exitoso
    			if t.Hour() < 0 || t.Hour() > 23 || t.Minute() < 0 || t.Minute() > 59 {
    				return nil, fmt.Errorf("hora o minutos fuera de rango")
    			}
    			return &t, nil
    		}
    	}
    
    	return nil, fmt.Errorf("formato de hora inválido, formatos aceptados: HH.MM, HH:MM, HHMM, HH MM")
    }
    ```
    
     
    
- **Función:** saveNotification
    
    **PROPOSITO:** recibe el struct con el modelo creado por **tranformRequestToModel** y guarda la entrada en la base de datos.
    
    **RETORNA:**
    
    1. Booleano.
    2. error.
    
    **PROCESO:**
    
    1. Revise que el modelo y el pointer de la BD no sea **nil**.
    2. Chequea que el modelo se cree exitosamente y que haya filas afectadas.
    
    ```go
    func saveNotification(model *models.NotificationBDV, db *gorm.DB) (bool, error) {
    	if model == nil {
    		return false, fmt.Errorf("notification model cannot be nil")
    	}
    
    	if db == nil {
    		return false, fmt.Errorf("database connection cannot be nil")
    	}
    
    	result := db.Create(model)
    	if result.Error != nil {
    		log.Printf("Error saving BDV notification: %v", result.Error)
    		return false, result.Error
    	}
    
    	if result.RowsAffected == 0 {
    		log.Println("No rows were affected when saving notification")
    		return false, fmt.Errorf("no rows affected")
    	}
    
    	return true, nil
    }
    
    ```
    
     
    
- **Función:** CheckNotificationExists
    
    **PROPOSITO:** revisa que la entrada no este duplicada usando los siguiente campos:
    
    - banco_origen
    - referencia_origen
    - fecha_banco
    - id_cliente
    
    **RETORNA:**
    
    1. Booleano.
    2. error.
    
    **PROCESO:**
    
    1. Revise que el modelo y el pointer de la BD no sea **nil**.
    2. Chequea que el modelo se cree exitosamente y que haya filas afectadas.
    
    ```go
    func CheckNotificationExists(bancoOrigen string, referenciaOrigen string, fechaBanco string,
    	id_cliente string, db *gorm.DB) (bool, error) {
    
    	if db == nil {
    		return false, fmt.Errorf("database connection cannot be nil")
    	}
    
    	var count int64
    
    	result := db.Model(&models.NotificationBDV{}).
    		Where("banco_origen = ? AND referencia_origen = ? AND fecha_banco = ? AND id_cliente = ?",
    			bancoOrigen, referenciaOrigen, fechaBanco, id_cliente).
    		Count(&count)
    
    	if result.Error != nil {
    		log.Printf("Error checking for existing BDV notification: %v", result.Error)
    		return false, result.Error
    	}
    
    	// Si count > 0, significa que ya existe al menos una notificación con esos datos
    	return count > 0, nil
    }
    
    ```
    
     
    
- **Función:** Validate
    
    **PROPOSITO:** Revisa que los campos de la petición no lleguen vacíos. 
    
    **RETORNA:**
    
    1. error.
    
    **PROCESO:**
    
    1. Si alg**ú**n campo se encuentra vac**í**o retorna el error correspondiente, de lo contrario el error se retorna como **nil.**
    
    ```go
    func (r *bdvRequest) Validate() error {
    	if r.BancoOrdenante == "" {
    		return fmt.Errorf("bancoOrdenante es obligatorio")
    	}
    	if r.Referencia == "" {
    		return fmt.Errorf("referenciaBancoOrdenante es obligatorio")
    	}
    	if r.IdCliente == "" {
    		return fmt.Errorf("idCliente es obligatorio")
    	}
    	if r.IdComercio == "" {
    		return fmt.Errorf("idComercio es obligatorio")
    	}
    	if r.NumeroCliente == "" {
    		return fmt.Errorf("numeroCliente es obligatorio")
    	}
    	if r.NumeroComercio == "" {
    		return fmt.Errorf("numeroComercio es obligatorio")
    	}
    	if r.Fecha == "" {
    		return fmt.Errorf("fecha es obligatorio")
    	}
    	if r.Hora == "" {
    		return fmt.Errorf("hora es obligatorio")
    	}
    	if r.Monto == "" {
    		return fmt.Errorf("monto es obligatorio")
    	}
    	return nil // Todos los campos están correctos
    }
    
    ```
    
     
    

---

---