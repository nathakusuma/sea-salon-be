Back-end repository: https://github.com/nathakusuma/sea-salon-be

Front-end repository:   https://github.com/nathakusuma/sea-salon-fe

# 💅 SEA SALON
Introducing SEA Salon, a distinguished name in the salon industry celebrated for its exceptional beauty services and stellar reviews. With a rapidly expanding clientele and a reputation for excellence, SEA Salon stands as the ultimate destination for all your beauty needs. Offering a comprehensive range of services, including haircuts and styling, manicures and pedicures, and facial treatments, the salon ensures a top-notch experience for every client. Featuring a user-friendly application, SEA Salon allows users to effortlessly explore services, choose their favorite stylists, and book appointments with ease. Embrace beauty and elegance redefined at SEA Salon, where your satisfaction is the priority.

## 📦 Deployment
Click the link below to access the website: </br>  
https://seasalon.nathakusuma.com

## 🛠️ Tech Stack
### Front-end
- React.js
- Bootstrap

### Back-end
- Golang
- PostgreSQL
- Fiber
- GORM
- Google Cloud Storage

## 🚀 How to run the code?
### Back-end
1. Make sure you have **Go** version **1.22** or later
2. Clone the github repository:
```
$ git clone https://github.com/nathakusuma/sea-salon-be.git
$ cd sea-salon-be
```
3. Provide all the environment variables in `.env` file. See `.env.example`.
4. Download the packages:
```
$ go mod tidy
```
5. Run the app:
```
$ go run cmd/app/main.go
```

This project uses Google Cloud Storage for storing assets. To run this website locally, you need to create a [bucket](https://console.cloud.google.com/storage). Make sure to not select/enable **Prevent Public Access**. After that, go to permissions tab and select **Grant Access**. Under **Add Principals**, type in **allUsers** and click **Save**. Finally, put the project id and bucket name in the environment variables.

Create a [service account](https://console.cloud.google.com/iam-admin/serviceaccounts). Under **Grant this service account access to project (optional)**, enter **Storage Object Admin** role. Click done. Then, select that service account, go to **Keys** tab, and click **Add Key > Create new key > JSON**. Convert the JSON file into **one line of JSON** and put it into the environment variable. If you're using **.env** file, then you need to **escape reserved characters** such as double quotes with backslash (example: \" becomes \\"). Make sure to escape the actual backslash (example: \\n becomes \\\\n).

### Front-end
1. Make sure you have **Node.js** version **v22.3.0** or later.
2. Clone the github repository:
```
$ git clone https://github.com/nathakusuma/sea-salon-fe.git
$ cd sea-salon-fe
```
2. Provide all the environment variables in `.env` file. See `.env.example`.
   Provide the backend URL as VITE_API_BASE_URL (example: http://localhost:8080/v1).

4. Download the packages:
```
$ npm install
```
5. Run the app:
```
$ npm run build
$ npm run preview
```