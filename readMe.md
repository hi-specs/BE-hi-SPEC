# hi'SPEC
<img src="https://github.com/hi-specs/FE-hi-SPEC/assets/127754894/8eaa9ccc-4e50-4066-a41f-bd3c481b404d">

<!-- ABOUT THE PROJECT -->
## About The Project
hi'Spec is a website that helps you choose the right laptop according to your needs. With a compare feature, you can freely bring the laptop side to side, and decide which will suitable for your needs.

## Build App & Database
![OpenAI Badge](https://img.shields.io/badge/OpenAI-412991?logo=openai&logoColor=fff&style=flat-square)
![Github Badge](https://img.shields.io/badge/Github-black?logo=github)
![Midtrans Badge](https://img.shields.io/badge/Midtrans-blue?logo=midtrans)
![Go Badge](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=fff&style=flat-square)
![Ubuntu Badge](https://img.shields.io/badge/Ubuntu-E95420?logo=ubuntu&logoColor=fff&style=flat-square)
![Cloudflare Badge](https://img.shields.io/badge/Cloudflare-F38020?logo=cloudflare&logoColor=fff&style=flat-square)
![Amazon EC2 Badge](https://img.shields.io/badge/Amazon%20EC2-F90?logo=amazonec2&logoColor=fff&style=flat-square)
![Swagger Badge](https://img.shields.io/badge/Swagger-85EA2D?logo=swagger&logoColor=000&style=flat-square)
![Postman Badge](https://img.shields.io/badge/Postman-FF6C37?logo=postman&logoColor=fff&style=flat-square)
![Insomnia Badge](https://img.shields.io/badge/Insomnia-4000BF?logo=insomnia&logoColor=fff&style=flat-square)
![Docker Badge](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=fff&style=flat-square)
![MySQL Badge](https://img.shields.io/badge/MySQL-4479A1?logo=mysql&logoColor=fff&style=flat-square)
![JSON Web Tokens Badge](https://img.shields.io/badge/JSON%20Web%20Tokens-000?logo=jsonwebtokens&logoColor=fff&style=flat-square)

## TechStack
![TectStack](https://github.com/hi-specs/BE-hi-SPEC/assets/73748420/e1769dd4-464b-4d56-9de3-1a65b598b495)

## ERD
![ERD_Hi'Spec (1)](https://github.com/hi-specs/BE-hi-SPEC/assets/50069221/e1302740-e0e9-49cc-b16d-c48698fac2a4)


## Run Locally

Clone the project

```bash
git clone https://github.com/hi-specs/BE-hi-SPEC.git
```

Go to the project directory

```bash
cd BE-hi-SPEC
```

Install dependency

```bash
go mod tidy
```

## Open Api 

If you're interested in using our Open Api, this is an example of how to do so.

Final Project Capstone Program Immersive Alterra Academy
<br />
<a href="https://app.swaggerhub.com/apis/hi_specs/hi_specs/1.0.0"><strong>Go to Open API »</strong></a>
<br />
<div>
      <details>
<summary>Admin</summary> 
<div>
  
| Feature User | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| GET | /dashboard  | - | + | Display the total user, product, and transaction. |
| GET | /users | - | + | Get all user. |
| GET | /transactions | - | + | Get all transaction. |
| POST | /product  | - | + | Create product with OpenAI. |

</details>

<div>
      <details>
<summary>User</summary>
  
  <!---
  | Command | Description |
| --- | --- |
  --->
 
<div>
  
| Feature Groups | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| GET | /user | - | + | Get user by id. |
| GET | /user/search | Email User | + | Search user. |
| POST | /login  | - | - | Login. |
| POST | /register | - | - | Register. |
| POST | /user/fav/add/{id} | Product ID | + | Wishlist. |
| PATCH | /user/{id} | User ID | + | Update user by id. |
| DELETE | /user/{id} | User ID | + | Displaying Group detail by id. |
| DELETE | /group/{id} | ID Groups | YES | Delete Groups. |
| DELETE | /user/fav/{id} | Wishlist ID | + | Delete Groups. |

</details>

<div>
      <details>
<summary>Product</summary>
  
  <!---
  | Command | Description |
| --- | --- |
  --->
 
<div>
  
| Feature Chats | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| GET | /products  | - | - | Get all product. |
| GET | /product/{id}  | Product ID | - | Get product. |
| GET | /product/search  | Name/Category/MinPrice/MaxPrice | - | Send a message to the groups. |
| PATCH | /product/{id}  | Product ID | + | Edit product. |
| DELETE | /product/{id}  | Product ID | + | Delete product. |

</details>

<div>
      <details>
<summary>Transaction</summary>
  
  <!---
  | Command | Description |
| --- | --- |
  --->
  
Several commands make use of Locations features, as shown below.
 
<div>
  
| Feature Locations | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| POST | /transaction  | - | + | Checkout product. |
| GET | /transaction/user/{id} | User ID | + | Get transaction list from user. |
| GET | /transaction/download/{id}  | Transaction ID | + | Download transaction detail. |

</details>

## Author

