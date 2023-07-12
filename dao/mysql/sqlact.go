package mysql

import (
	"bcloud/netdisk/download"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

func InsertListDownDetail(db *sql.DB, list []*download.DownDetail) {
	query := "INSERT INTO tb_down_detail (created_at, name,path,size,status,dlink,fsid,process_status) VALUES (?,?,?,?,?,?,?,?)"
	prepare, _ := db.Prepare(query)
	for _, i := range list {
		if _, err := prepare.Exec(i.CreatedAt, i.Name, i.Path, i.Size, i.Status, i.Dlink, i.Fsid, 0); err != nil {
			zap.L().Error("下载列表插入数据库失败", zap.Error(err))
		}
	}
	zap.L().Info("下载列表插入数据库成功")
}

func FindListDownDetail(db *sql.DB) []*download.DownDetail {
	rows, err := db.Query("SELECT id, created_at, name, path, size, status, dlink, fsid, process_status FROM tb_down_detail WHERE status <> ?", "下载完成")
	if err != nil {
		zap.L().Error("查询数据失败", zap.Error(err))
	}
	defer rows.Close()
	var downDetails []*download.DownDetail
	for rows.Next() {
		var d download.DownDetail
		err := rows.Scan(
			&d.ID,
			&d.CreatedAt,
			&d.Name,
			&d.Path,
			&d.Size,
			&d.Status,
			&d.Dlink,
			&d.Fsid,
			&d.ProcessStatus,
		)
		if err != nil {
			fmt.Println("扫描行数据失败:", err)
			return nil
		}
		downDetails = append(downDetails, &d)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("遍历结果集失败:", zap.Error(err))
		return nil
	}
	return downDetails
}

func FindListOrderId(db *sql.DB) []download.DownDetail {
	rows, _ := db.Query("SELECT id, created_at, name, path, size, status, dlink, fsid, process_status FROM tb_down_detail ORDER BY id DESC")
	defer rows.Close()
	var downDetails []download.DownDetail
	for rows.Next() {
		var d download.DownDetail
		_ = rows.Scan(
			&d.ID,
			&d.CreatedAt,
			&d.Name,
			&d.Path,
			&d.Size,
			&d.Status,
			&d.Dlink,
			&d.Fsid,
			&d.ProcessStatus,
		)
		downDetails = append(downDetails, d)
	}
	return downDetails
}

func FindDownDetailForId(db *sql.DB, id int) download.DownDetail {
	var d download.DownDetail
	_ = db.QueryRow("SELECT id, created_at, name, path, size, status, dlink, fsid, process_status FROM tb_down_detail WHERE id = ?", id).Scan(&d.ID,
		&d.CreatedAt,
		&d.Name,
		&d.Path,
		&d.Size,
		&d.Status,
		&d.Dlink,
		&d.Fsid,
		&d.ProcessStatus)
	return d
}
