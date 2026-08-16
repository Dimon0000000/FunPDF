// Last modified: 2026-08-16
//
// Copyright 2026 Haruko386, SJZU. All rights reserved.

//  Licensed under the GNU GENERAL PUBLIC LICENSE, Version 3.0 (the "License");
//  you may not use this file except in compliance with the License.
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package main

import (
	"FunPDF/internal"
	"FunPDF/internal/common"
	"FunPDF/internal/dao"
	"FunPDF/internal/entity"
	"FunPDF/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fileHandler := handler.NewFileHandler()

	r := gin.Default()
	router := internal.NewRouter(fileHandler)
	router.Setup(r)

	dsn := "root:password@(127.0.0.1:3306)/funpdf?charset=utf8&parseTime=True&loc=Local"
	if err := dao.InitMysql(dsn); err != nil {
		log.Fatalf("initialize MySQL: %v", err)
	}
	if err := dao.DB.AutoMigrate(&entity.File{}); err != nil {
		log.Fatalf("migrate file table: %v", err)
	}

	common.Banner()
	log.Printf("FunPDF %s", common.GetVersion())

	if err := r.Run(":9384"); err != nil {
		log.Fatalf("start backend: %v", err)
	}
}
