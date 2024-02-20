// 登录日志 数据层
// author xiaoRui

package dao

import (
	"admin-api/api/entity"
	"admin-api/common/util"
	. "admin-api/pkg/db"
	"time"
)

// CreateSysLoginInfo 新增登录日志
func CreateSysLoginInfo(username, ipAddress, loginLocation, browser, os, message string, loginStatus int) {
	sysLoginInfo := entity.SysLoginInfo{
		Username:      username,
		IpAddress:     ipAddress,
		LoginLocation: loginLocation,
		Browser:       browser,
		Os:            os,
		Message:       message,
		LoginStatus:   loginStatus,
		LoginTime:     util.HTime{Time: time.Now()},
	}
	Db.Save(&sysLoginInfo)
}

// GetSysLoginInfoList 分页获取登录日志列表
func GetSysLoginInfoList(Username, LoginStatus, BeginTime, EndTime string, PageSize, PageNum int) (sysLoginInfo []entity.SysLoginInfo, count int64) {
	curDb := Db.Table("sys_login_info")
	if Username != "" {
		curDb = curDb.Where("username = ?", Username)
	}
	if BeginTime != "" && EndTime != "" {
		curDb = curDb.Where("`login_time` BETWEEN ? AND ?", BeginTime, EndTime)
	}
	if LoginStatus != "" {
		curDb = curDb.Where("login_status = ?", LoginStatus)
	}
	curDb.Count(&count)
	curDb.Limit(PageSize).Offset((PageNum - 1) * PageSize).Order("login_time desc").Find(&sysLoginInfo)
	return sysLoginInfo, count
}

// BatchDeleteSysLoginInfo 批量删除登录日志
func BatchDeleteSysLoginInfo(dto entity.DelSysLoginInfoDto) {
	Db.Where("id in (?)", dto.Ids).Delete(&entity.SysLoginInfo{})
}

// DeleteSysLoginInfoById 根据id删除日志
func DeleteSysLoginInfoById(dto entity.SysLoginInfoIdDto) {
	Db.Delete(&entity.SysLoginInfo{}, dto.Id)
}

// CleanSysLoginInfo 清空登录日志
func CleanSysLoginInfo() {
	Db.Exec("truncate table sys_login_Info")
}
