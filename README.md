# Akali 專案部署說明

## 專案概述
Akali 為 Kubernetes (GKE) 上運行的服務，主要結合 Cloud SQL 與外部 API，並提供健康檢查機制與後台管理介面。

---

## 基本資訊

| 項目 | 值 |
|------|----|
| **GKE Cluster 名稱** | `????` |
| **GCP 區域** | `????` |
| **Namespace** | `akali` |
| **服務訪問 IP** | [http://localhost/](http://localhost/) |
| **健康檢查路徑** | [http://localhost/healthy](http://localhost/healthy) |

---

## Cloud SQL 設定

| 項目 | 值 |
|------|----|
| **DB Host (Private IP)** | `XXX.XXX.XXX.XXX` |
| **DB 使用者** | `root` |
| **DB 密碼** | `??????` |
| **連線方式** | Private IP 連線 (同 VPC 內) |

---

## phpMyAdmin 操作說明

可透過 Port Forward 方式連線至 Kubernetes 內部的 phpMyAdmin 服務。

```bash
