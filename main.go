package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Request struct {
	Query    string `json:"query"`
	ApiKey   string `json:"apikey"`
	Values   []any  `json:"values"`
	Database string `json:"db"`
}

var (
	databases = make(map[string]*badger.DB)
	dbMutex   = sync.RWMutex{}
)

func main() {
	conf, err := os.ReadFile("dbconfig.json")
	var config map[string]string
	if err != nil {
		config = map[string]string{
			"apikey": "gorms",
			"dbport": "3308",
		}
		jsonData, _ := json.MarshalIndent(config, "", " ")
		os.WriteFile("dbconfig.json", jsonData, 0777)
	} else {
		json.Unmarshal(conf, &config)
	}
	jjj := bytes.NewReader([]byte("hola"))
	_, err = http.Post("http://127.0.0.1:"+config["dbport"]+"/", "application/json", jjj)
	if err == nil {
		return
	}
	/*db, err := OpenBadgerDB(config["directory"])
	if err != nil {
		panic(err)
	}
	kkdhooo := uuid.New().String()
	InsertData(db, kkdhooo, []any{kkdhooo})*/
	/*go func() {
		ticker := time.NewTicker(60 * time.Minute)
		for range ticker.C {
			db.RunValueLogGC(0.7)
		}
	}()*/
	os.MkdirAll("backups", 0755)
	os.MkdirAll("databases", 0755)
	r := GinRouter()
	fmt.Println("server ir running...", config["dbport"])
	r.POST("/data", func(ctx *gin.Context) {
		ex(ctx)
	})
	r.POST("/", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.Run("0.0.0.0:" + config["dbport"])
}

func GetOrCreateDB(dbName string) (*badger.DB, error) {
	dbMutex.RLock()
	if db, exists := databases[dbName]; exists {
		dbMutex.RUnlock()
		return db, nil
	}
	dbMutex.RUnlock()

	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Double-check después del lock
	if db, exists := databases[dbName]; exists {
		return db, nil
	}

	// Crear nueva base de datos
	dbPath := "./databases/" + dbName
	db, err := OpenBadgerDB(dbPath)
	if err != nil {
		return nil, err
	}

	databases[dbName] = db
	return db, nil
}

// Cerrar una base de datos específica
func CloseDB(dbName string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db, exists := databases[dbName]; exists {
		err := db.Close()
		delete(databases, dbName)
		return err
	}
	return nil
}

// Cerrar todas las bases de datos
func CloseAllDBs() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	for name, db := range databases {
		db.Close()
		delete(databases, name)
	}
}

func configurate() map[string]string {
	conf, err := os.ReadFile("dbconfig.json")
	var config map[string]string
	if err != nil {
		config = map[string]string{
			"apikey": "gorms",
			"dbport": "3308",
		}
		jsonData, _ := json.MarshalIndent(config, "", " ")
		os.WriteFile("dbconfig.json", jsonData, 0777)
		return config
	} else {
		json.Unmarshal(conf, &config)
		return config
	}
}

func ex(ctx *gin.Context) {
	var mapa Request
	err := ctx.ShouldBindJSON(&mapa)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
		return
	}
	if len(mapa.Query) > 1000 {
		ctx.JSON(400, gin.H{"status": "error", "message": "query too long", "result": []any{}})
		return
	}
	mapa.Query = strings.TrimSpace(mapa.Query)
	comando := strings.Split(mapa.Query, " ")
	fmt.Println(comando)
	//fmt.Println(mapa.Values...)
	if len(comando) < 2 {
		ctx.JSON(400, gin.H{"status": "error", "message": "invalid sintax", "result": []any{}})
		return
	}
	if len(comando) > 2 {
		comando[1] = strings.Join(comando[1:], " ")
	}
	comando1 := strings.ToUpper(comando[0])
	comando2 := comando[1]
	conf := configurate()
	if conf["apikey"] != mapa.ApiKey {
		ctx.JSON(401, gin.H{"status": "error", "message": "invalid apikey", "result": []any{}})
		return
	}
	if mapa.Database == "" {
		ctx.JSON(400, gin.H{"status": "error", "message": "database name empty", "result": []any{}})
		return
	}
	db, err := GetOrCreateDB(mapa.Database)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
		return
	}
	if comando1 == "BACKUP" {
		err = CreateBackup(db, "./backups/"+mapa.Database+"_"+GetFormattedDate()+".bak")
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "backup created", "result": []any{}})
		return
	}
	if comando1 == "LIKE" {
		dat, err := QueryByPrefix(db, comando2, mapa.Values)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "result": dat})
		return
	}
	if comando1 == "SELECT" {
		dat, code, err := GetData(db, comando2)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(code, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "result": dat})
		return
	}
	if comando1 == "COMPRESS" {
		err := db.RunValueLogGC(0.7)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "result": []any{}})
		return
	}

	if comando1 == "DELETE" {
		err := DeleteData(db, comando2)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "result": []any{}})
		return
	}
	if comando1 == "UPDATE" {
		err := UpdateData(db, comando2, mapa.Values)
		if err != nil {
			fmt.Println(err)
			code := 500
			if err.Error() == "key does not exist" {
				code = 404
			}
			ctx.JSON(code, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "result": []any{}})
		return
	}
	if comando1 == "INSERT" {
		err := InsertData(db, comando2, mapa.Values)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "result": []any{}})
		return
	}
	if comando1 == "RESTORE" {
		err = RestoreBackup(db, "./"+comando2)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "result": []any{}})
			return
		}
		ctx.JSON(200, gin.H{"status": "success", "message": "database restored", "result": []any{}})
		return
	}
	ctx.JSON(400, gin.H{"status": "error", "message": "invalid sintax", "result": []any{}})
}

func GetFormattedDate() string {
	return time.Now().Format("2006-01-02_15_04_05")
}

func RestoreBackup(db *badger.DB, backupPath string) error {
	file, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return db.Load(file, 16)
}

// Obtener estadísticas de una base de datos
func GetDBStats(dbName string) (map[string]interface{}, error) {
	dbMutex.RLock()
	db, exists := databases[dbName]
	dbMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("database not found")
	}

	lsm, vlog := db.Size()
	return map[string]interface{}{
		"lsm_size":  lsm,
		"vlog_size": vlog,
		"name":      dbName,
	}, nil
}

// Backup de una base de datos específica
func BackupSpecificDB(dbName, backupPath string) error {
	dbMutex.RLock()
	db, exists := databases[dbName]
	dbMutex.RUnlock()

	if !exists {
		return fmt.Errorf("database not found")
	}

	return CreateBackup(db, backupPath)
}

func GetData(db *badger.DB, key string) ([]any, int, error) {
	var value []byte
	result := make([]any, 0)
	code := 0
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			code = 404
			return err
		}

		value, err = item.ValueCopy(nil)
		code = 500
		return err
	})

	if err != nil {
		return nil, code, err
	}
	err = json.Unmarshal(value, &result)
	if err != nil {
		return nil, 500, err
	}

	return result, 200, err
}

func QueryByPrefix(db *badger.DB, prefix string, numero []any) (map[string][]any, error) {
	result := make(map[string][]any)
	nu := numero[0].(float64)
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefixBytes := []byte(prefix)
		var cont int64 = 0
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			if cont >= int64(nu) {
				break
			}
			item := it.Item()
			key := item.Key()
			val, err := item.ValueCopy(nil)
			hhsb := make([]any, 0)
			if err != nil {
				continue
			}
			err = json.Unmarshal(val, &hhsb)
			if err != nil {
				continue
			}
			result[string(key)] = hhsb
			cont++

		}
		return nil
	})

	return result, err
}

func DeleteData(db *badger.DB, key string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func UpdateData(db *badger.DB, key string, value []any) error {
	return db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return fmt.Errorf("key does not exist")
		}
		if err != nil {
			return err
		}

		valueBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), valueBytes)
	})
}

func OpenBadgerDB(path string) (*badger.DB, error) {
	opts := badger.DefaultOptions(path)
	opts.SyncWrites = true
	opts.DetectConflicts = false
	opts.CompactL0OnClose = false
	opts.MemTableSize = 16 << 20
	opts.NumMemtables = 2
	opts.NumLevelZeroTables = 2
	opts.Logger = nil

	// Intento 1: Apertura normal
	db, err := badger.Open(opts)
	if err == nil {
		return db, nil
	}
	fmt.Printf("Normal open failed: %v\n", err)

	// Intento 2: Bypass lock guard
	opts.BypassLockGuard = true
	db, err = badger.Open(opts)
	if err == nil {
		fmt.Println("Recovered using BypassLockGuard")
		return db, nil
	}
	fmt.Printf("BypassLockGuard failed: %v\n", err)

	// Intento 3: Modo solo lectura
	opts.ReadOnly = true
	db, err = badger.Open(opts)
	if err == nil {
		fmt.Println("Opened in ReadOnly mode")
		return db, nil
	}
	fmt.Printf("ReadOnly mode failed: %v\n", err)

	// Intento 4: Configuración mínima
	opts = badger.DefaultOptions(path)
	opts.MemTableSize = 1 << 20     // 1MB
	opts.ValueLogFileSize = 1 << 20 // 1MB - COMPATIBLE v4
	opts.NumMemtables = 1
	opts.NumLevelZeroTables = 1
	opts.SyncWrites = false // Más permisivo
	opts.BypassLockGuard = true
	opts.Logger = nil

	db, err = badger.Open(opts)
	if err == nil {
		fmt.Println("Recovered with minimal settings")
		return db, nil
	}

	// Intento 5: Backup y recrear (sin cambios)
	backupPath := path + "_backup_" + fmt.Sprintf("%d", time.Now().Unix())
	if err := os.Rename(path, backupPath); err == nil {
		fmt.Printf("Original DB moved to: %s\n", backupPath)

		opts = badger.DefaultOptions(path)
		opts.SyncWrites = true
		db, err = badger.Open(opts)
		if err == nil {
			fmt.Printf("New DB created. Backup: %s\n", backupPath)
			return db, nil
		}
	}

	return nil, fmt.Errorf("all recovery attempts failed: %v", err)
}
func CreateBackup(db *badger.DB, backupPath string) error {
	file, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = db.Backup(file, 0)
	return err
}

func InsertData(db *badger.DB, key string, value []any) error {
	return db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == nil {
			fmt.Println(key)
			return fmt.Errorf("key already exists")
		}
		if err != badger.ErrKeyNotFound {
			return err
		}

		valueBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), valueBytes)
	})
}

func SetupLogger() {
	// Crear o abrir el archivo de logs
	file, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Print(err)
	}

	// Configurar el logger para escribir en el archivo
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func errorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Procesar la solicitud
		c.Next()

		// Comprobar si hubo errores
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Println(err.Error())
			}
		}
	}
}

func GinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(errorLogger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	return r
}
