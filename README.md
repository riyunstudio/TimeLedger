# Akali 專案部署說明

## 專案概述
Akali 為 Kubernetes (GKE) 上運行的服務，主要結合 Cloud SQL 與外部 API，並提供健康檢查機制與後台管理介面。

---

## 基本資訊

| 項目 | 值 |
|------|----|
| **GKE Cluster 名稱** | `bots` |
| **GCP 區域** | `asia-east1-b` |
| **Namespace** | `akali` |
| **服務訪問 IP** | [http://34.81.134.97/](http://34.81.134.97/) |
| **健康檢查路徑** | [http://34.81.134.97/healthy](http://34.81.134.97/healthy) |

---

## Cloud SQL 設定

| 項目 | 值 |
|------|----|
| **DB Host (Private IP)** | `10.140.0.3` |
| **DB 使用者** | `master` |
| **DB 密碼** | `=H[7u3Va;fl2qGo[` |
| **連線方式** | Private IP 連線 (同 VPC 內) |

---

## phpMyAdmin 操作說明

可透過 Port Forward 方式連線至 Kubernetes 內部的 phpMyAdmin 服務。

```bash
gcloud container clusters get-credentials bots --zone asia-east1-b --project telegram-bet \
 && kubectl port-forward --namespace phpmyadmin $(kubectl get pod --namespace phpmyadmin --selector="app=phpmyadmin-master" --output jsonpath='{.items[0].metadata.name}') 8080:80