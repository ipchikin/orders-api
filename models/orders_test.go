package models

// // TestConnect to db
// func TestConnect(t *testing.T) {
// 	cfg, err := configs.LoadConfig("test")
// 	assert.Nil(t, err)

// 	bm := new(BaseModel)
// 	err = bm.Connect(
// 		cfg.DBConfig.Driver,
// 		cfg.DBConfig.User,
// 		cfg.DBConfig.Password,
// 		cfg.DBConfig.Host,
// 		cfg.DBConfig.Port,
// 		cfg.DBConfig.Database,
// 		cfg.DBConfig.MaxIdleConns,
// 	)
// 	assert.Nil(t, err)
// }

// // TestConnectErr with invalid db config
// func TestConnectErr(t *testing.T) {
// 	bm := new(BaseModel)
// 	err := bm.Connect("", "", "", "", "", "", 0)
// 	assert.NotNil(t, err)
// }
