# Các chạy dự án local

## TODO

- docker + docker compose cho dev environment

## Các chạy local

- cài đặt đủ các package trong requirement.txt
- cài đặt foundry [link](https://book.getfoundry.sh/)

- chạy local evm node
```
cd smartcontract
make anvil-up
```
- local env node phải có file smartcontract/cache/anvil_state.json (có thể tải trong [link](https://drive.google.com/drive/folders/1gdCsQFRMT3Egs7zYiCoCXARmL_QHTsFu?usp=drive_link]))
- nếu chưa có file anvil_state cần deploy lại contract 
    ```
        // in smartcontract folder 
        make deploy
    ```
- nếu deploy lại contract, chú ý cập nhật lại contract address trong server/config/config.yaml và web/.env

- chạy backend
```
cd server
go run cmd/v2/main.go run
```

- chay frontend
```
cd server
go run cmd/v2/main.go run
```