# Hito 1 · Ficha del negocio

**Pareja:** Cedeño Francisco · Santana Alexis
**Paralelo:**  Aplicaciones Web II A 
**Negocio en una línea:** Aplicación web para dueños y veterinarios que gestiona mascotas y sus registros de salud y actividades.



## 1. Negocio de referencia

El negocio de referencia es PetDesk, una aplicación SaaS que conecta clínicas veterinarias con dueños de mascotas para gestionar citas, recordatorios de vacunas y el historial médico digital. Vende servicios de gestión de salud y actividades, cobra mediante suscripción mensual a las clínicas y declara ingresos superiores a 5 millones de dólares anuales. Su éxito se debe a que resuelve un problema real: la falta de organización en citas y tratamientos, que genera pérdidas de tiempo y olvidos en la atención. Los fundadores explican que la clave fue integrar recordatorios automáticos y comunicación directa con las clínicas veterinarias, lo que generó fidelización y uso recurrente. Además, la aplicación permite a los dueños tener toda la información de sus mascotas centralizada en un solo lugar, lo que aumenta la confianza y la satisfacción del cliente. Este caso demuestra que la gestión digital de mascotas es un nicho viable y escalable, especialmente cuando se combina con un modelo de ingresos sostenibles y un mercado amplio de dueños que buscan soluciones prácticas.

**Enlace:** https://www.capterra.ec/software/166086/petdesk

## 2. Caso de contraste

El caso de contraste es VetPronto, una startup estadounidense que ofrecía veterinarios a domicilio mediante una aplicación móvil. Fundada en 2014, buscaba facilitar la atención médica de mascotas en casa, pero cerró operaciones en 2018. La hipótesis principal de su fracaso es que el modelo dependía de una logística compleja y costosa: coordinar visitas presenciales implicaba altos gastos operativos y baja escalabilidad. Además, la adopción digital fue limitada, ya que muchos dueños preferían acudir directamente a clínicas veterinarias por confianza y disponibilidad de equipos médicos. VetPronto no logró alcanzar masa crítica de usuarios y sus costos superaron los ingresos, lo que llevó al cierre. La diferencia de fondo con nuestro negocio es que VetPronto intentó resolver el problema con desplazamientos físicos, mientras que nuestra aplicación se centra en la gestión digital de registros y comunicación directa, lo que reduce costos y asegura un modelo más sostenible y escalable.

**Fuente:** https://medium.com/@joewaltman/vetpronto-is-being-euthanized-a2aa39fbd78f

## 3. Adaptación al Ecuador


1. Medios de pago limitados: En Ecuador muchos clientes prefieren pagar mediante transferencia bancaria o efectivo, enviando comprobantes por WhatsApp. Esto obliga a incluir un estado adicional en el flujo de pagos: “pendiente de verificación”, antes de activar la suscripción o el servicio.
2. Conectividad irregular fuera de ciudades grandes: En zonas rurales la conexión a internet es irregular. La aplicación debe permitir registrar datos de mascotas en modo offline y sincronizar la información cuando haya acceso a la red, para no perder registros importantes.
3. Normativa tributaria del SRI: La facturación electrónica es obligatoria y los negocios deben registrarse con RUC o RIMPE. Esto implica que el sistema debe generar comprobantes electrónicos y adaptarse a los formatos exigidos por el SRI, asegurando cumplimiento legal.

**Qué cambió en el modelo por estas restricciones:** 

Se agregó un estado de pago “pendiente de verificación” en la máquina de estados, se diseñó un módulo de sincronización offline para zonas con baja conectividad y se integró la facturación electrónica en el flujo de cobros. Estos cambios garantizan que el negocio sea viable en el contexto ecuatoriano, reduciendo fricciones y cumpliendo con la normativa local.

## 4. Modelo de datos


### Entidad: Mascota

| Atributo  |       Tipo        | Obligatorio |   Ejemplo        |
|---------- |-------------------|-------------|------------------|
|    ID     |  numero entero    |      si     |       1          |
|   NOMBRE  |    texto          |      si     |       Firulais   |
|   ESPECIE |    texto          |      si     |       Perro      |
|   ESTADO  |    texto          |      si     |       Saludable  |
|   FechaNac|    fecha          |      si     |       05/12/2020 |
|   DueñoID |    numero entero  |      si     |       5          |

### Entidad: Registro

| Atributo  |       Tipo        | Obligatorio |   Ejemplo              |
|---------- |-------------------|-------------|------------------------|
|    ID     |  numero entero    |      si     |       100              |
| MascotaID |  numero entero    |      si     |       10               |
|   Tipo    |    texto          |      si     |       vacuna           |
|   Detalle |    texto          |      no     |       rabia            |
|   Fecha   |    fecha y hora   |      si     |       05/12/2026 10:00 |

### Entidad: Usuario

| Atributo  |       Tipo        | Obligatorio |   Ejemplo             |
|---------- |-------------------|-------------|-----------------------|
|    ID     |  numero entero    |      si     |       1               |
|   NOMBRE  |    texto          |      si     |       Francisco       |
|   rol     |    texto          |      si     |       dueño           |
|   correo  |    texto          |      si     |  francsico@gmail.com  |

### Entidad: Cita

| Atributo  |       Tipo        | Obligatorio |   Ejemplo              |
|---------- |-------------------|-------------|------------------------|
|    ID     |  numero entero    |      si     |       105              |
| MascotaID |  numero entero    |      si     |       19               |
|   Motivo  |    texto          |      si     |       chequeo          |
|   Fecha   |    fecha y hora   |      si     |       05/12/2026 10:00 |

### Entidad: Tratamiento

| Atributo  |       Tipo        | Obligatorio |   Ejemplo              |
|---------- |-------------------|-------------|------------------------|
|    ID     |  numero entero    |      si     |       100              |
| MascotaID |  numero entero    |      si     |       10               |
|descripcion|    texto          |      si     | antibióticos por 5 dias|
|   iniciotr|    fecha y hora   |      si     |       28/11/2026 09:30 |
|   Fintr   |    fecha y hora   |      no     |       05/12/2026 10:00 |

### Entidad: Pago

| Atributo    |       Tipo        | Obligatorio |   Ejemplo        |
|-------------|-------------------|-------------|------------------|
|    ID       |  numero entero    |      si     |       1          |
|   usuarioID |  numero entero    |      si     |       100        |
|   monto     |  numero decimal   |      si     |       29.99      |
|   ESTADO    |    texto          |      si     |       verificado |
|   Fecha     |    fecha          |      si     |       05/12/2020 |
| comprobante |    texto          |      si     |       FAC-0001   |

### Relaciones


|      Entidades      | Cardinalidad    |                                     Frase                                               |
|---------------------|-----------------|-----------------------------------------------------------------------------------------|
| Usuario-Mascota     |      1:N        | Un usuario puede registrar varias mascotas, cada mascota pertenece a un usuario.        |
| Mascota-Registro    |      1:N        | Una mascota tiene muchos registros, cada registro pertenece a una sola mascota.         |
| Mascota-Cita        |      1:N        | Una mascota puede tener varias citas, cada cita corresponde a una mascota.              |
| Mascota-Tratamiento |      1:N        | Una mascota puede tener varios tratamientos, cada tratamiento pertenece a una mascota.  |
| Usuario-Pago        |      1:N        | Un usuario puede generar varias facturas/pagos, cada factura pertenece a un usuario.    |


### Structs en Go

type Usuario struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    Nombre string `gorm:"not null" json:"nombre"`
    Rol    string `gorm:"not null" json:"rol"` // dueño o veterinario
    Correo string `gorm:"unique;not null" json:"correo"`
}

type Mascota struct {
    ID       uint      `gorm:"primaryKey" json:"id"`
    Nombre   string    `gorm:"not null" json:"nombre"`
    Especie  string    `gorm:"not null" json:"especie"`
    Estado   string    `gorm:"not null;default:saludable" json:"estado"`
    FechaNac time.Time `json:"fecha_nac"`
    DuenoID  uint      `gorm:"not null" json:"dueno_id"`
}

type Registro struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    MascotaID uint      `gorm:"not null" json:"mascota_id"`
    Tipo      string    `gorm:"not null" json:"tipo"`
    Detalle   string    `json:"detalle"`
    Fecha     time.Time `gorm:"not null" json:"fecha"`
}

type Cita struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    MascotaID uint      `gorm:"not null" json:"mascota_id"`
    FechaHora time.Time `gorm:"not null" json:"fecha_hora"`
    Motivo    string    `gorm:"not null" json:"motivo"`
}

type Tratamiento struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    MascotaID uint      `gorm:"not null" json:"mascota_id"`
    Descripcion string  `gorm:"not null" json:"descripcion"`
    Inicio    time.Time `gorm:"not null" json:"inicio"`
    Fin       time.Time `json:"fin"`
}

type Pago struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UsuarioID   uint      `gorm:"not null" json:"usuario_id"`
    Monto       float64   `gorm:"not null" json:"monto"`
    Fecha       time.Time `gorm:"not null" json:"fecha"`
    Estado      string    `gorm:"not null;default:pendiente" json:"estado"`
    Comprobante string    `gorm:"not null" json:"comprobante"`
}



**Decisión de tipos que tuvimos que pensar:** 

-Para las fechas y horas usamos time.Time en lugar de texto, porque permite validar automáticamente el formato, realizar operaciones de comparación y evitar errores de escritura.

-Para el monto en Pago usamos float64 en lugar de número entero, ya que los valores monetarios pueden incluir decimales.

-Para los roles de Usuario y los estados de Mascota y Pago usamos cadenas (string) con listas cerradas (dueño/veterinario, saludable/en_observacion/enferma, pendiente/verificado/pagado). Esto asegura flexibilidad pero mantiene control sobre los valores permitidos.

-Para los IDs de todas las entidades usamos uint (número entero sin signo), porque nunca serán negativos y es el estándar en GORM para claves primarias.

-En el campo Fin de Tratamiento lo dejamos opcional (time.Time sin restricción de not null), porque algunos tratamientos pueden no tener fecha de finalización definida al inicio.

### Diagrama del modelo completo

erDiagram
    USUARIO ||--o{ MASCOTA : "registra"
    MASCOTA ||--o{ REGISTRO : "tiene"
    MASCOTA ||--o{ CITA : "agenda"
    MASCOTA ||--o{ TRATAMIENTO : "recibe"
    USUARIO ||--o{ PAGO : "realiza"

    USUARIO {
        int id
        string nombre
        string rol
        string correo
    }

    MASCOTA {
        int id
        string nombre
        string especie
        string estado
        date fecha_nac
        int dueno_id
    }

    REGISTRO {
        int id
        int mascota_id
        string tipo
        string detalle
        datetime fecha
    }

    CITA {
        int id
        int mascota_id
        datetime fecha_hora
        string motivo
    }

    TRATAMIENTO {
        int id
        int mascota_id
        string descripcion
        date inicio
        date fin
    }

    PAGO {
        int id
        int usuario_id
        float monto
        date fecha
        string estado
        string comprobante
    }


**Decisión discutible del modelo y por qué la tomamos:**

Se decidio quitar el campo veterinarioID de Cita ya que se esta usando en el rol del usuario pero en un ambiente real es necesario saber que veterinario atiende a que mascota

## 5. Máquina de estados

**Entidad con estados:** 

|       Estado       |                        Qué significa                             |
|--------------------|------------------------------------------------------------------|
| (inicial) saludable| La mascota está en buen estado de salud.                         |
| en_observacion     | La mascota presenta síntomas leves o está bajo control.          |
| enfermo            | La mascota requiere atención médica activa.                      |
| recuperado         | La mascota ha superado un tratamiento y vuelve a estar estable.  |


| De             | A               | Quién la hace |                           Condición                                |
|----------------|-----------------|---------------|--------------------------------------------------------------------|
| saludable      | en_observacion  |  Veterinario  | Se detectan síntomas leves en una cita o registro.                 |
| en_observacion | enfermo         |  Veterinario  | Se confirma diagnóstico de enfermedad.                             |
| enfermo        | recuperado      |  Veterinario  | Tratamiento finalizado con éxito.                                  |
| recuperado     | saludable       |  Veterinario  | Se valida que la mascota está estable en controles posteriores.    |

**Transición prohibida y por qué:** 

De saludable => recuperado: No tiene sentido que una mascota “recupere” sin haber estado enferma o en observación primero. La recuperación solo aplica después de un tratamiento.

### Diagrama de estados

stateDiagram-v2
    [*] --> saludable
    saludable --> en_observacion : síntomas leves detectados
    en_observacion --> enferma : diagnóstico confirmado
    enferma --> recuperada : tratamiento exitoso
    recuperada --> saludable : control posterior

    note right of recuperada
        De saludable no se puede pasar directamente a recuperada:
        la recuperación solo aplica después de enfermedad u observación.
    end note


## 6. Roles y permisos



|        Acción            |      Dueño     | Veterinario |
|--------------------------|----------------|-------------|
| Registrar mascota        | sí             | no          |
| Consultar mascotas       | solo las suyas | todos       |
| Crear registro de salud  | no             | sí          |
| Consultar registros      | solo los suyos | todos       |
| Agendar cita             | sí             | sí          |
| Registrar tratamiento    | no             | sí          |
| Consultar tratamientos   | solo los suyos | todos       |
| Realizar pago            | sí             | no          |
| Consultar pagos          | solo los suyos | todos       |

## 7. Mapa de endpoints por rol


| Endpoint           | Rol que lo llama   | Pantalla que lo consume |           Qué devuelve        |               Qué valida              | Código si falla |
|--------------------|--------------------|-------------------------|-------------------------------|---------------------------------------|-----------------|
| GET /mascotas      | Dueño, veterinario | Lista de mascotas       | JSON con mascotas             | Autenticación y rol                   | 401 / 403       |
| POST /mascotas     | Dueño              | Registro de mascota     | Mascota creada                | Datos obligatorios (nombre, especie)  | 400             |
| GET /mascotas/{id} | Dueño, veterinario | Detalle de mascota      | JSON con datos de la mascota  | Existencia de la mascota              | 404             |
| PUT /mascotas/{id} | Dueño              | Edición de mascota      | Mascota actualizada           | Dueño es propietario                  | 403             |
| POST /citas        | Dueño              | Agenda de citas         | Cita creada                   | Fecha válida, mascota existente       | 400             |
| POST /pagos        | Dueño              | Pantalla de pagos       | Pago registrado               | Monto > 0, comprobante único          | 400             |

### Matriz pantalla × endpoint


| Pantalla             | [GET /mascotas]| [POST /mascotas]| [GET /mascotas/{id}]         | [PUT /mascotas/{id}]| [POST /mascotas] | [POST /pagos] |
|----------------------|----------------|-----------------|------------------------------|---------------------|------------------|---------------|
| Lista de mascotas    |       x        |                 |                              |                     |                  |               |
| Registro de mascotas |                |        x        |                              |                     |                  |               |
| Detalle de mascota   |                |                 |             x                |          x          |                  |               |
| Agenda de citas      |                |                 |                              |                     |         x        |               |
| Pantalla de pagos    |                |                 |                              |                     |                  |       x       |


**Endpoints que ya están funcionando y en qué archivo:**
 
GET /mascotas => controllers/mascota.go

POST /mascotas => controllers/mascota.go

GET /mascotas/{id} => controllers/mascota.go

PUT /mascotas/{id} => controllers/mascota.go

POST /citas => controllers/cita.go

POST /pagos => controllers/pago.go

## 8. Declaración de IA

Usamos Microsoft Copilot (IA generativa) para apoyar en la redacción y estructuración de las secciones del modelo