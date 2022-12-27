// let namaVariable = "isivariabel"

// const angka = 5

// const isMarried = true

// type data non primitive
/*
Array adalah tipe data non primitive yang menggunakan Indexing sebagai pengurut
pada tiap index array, dapat menggunakan tipe data yang berbeda
*/
// Array
// let array = ["Wahyu", 17, false]

// console.log(array)
// console.table(array[0])

// Object
// let myBio = {
//   name: "Dandi Saputra",
//   age: 17,
//   isMarried: 
//   when:{
//     age2: 24
//   }
// }

// console.table(`Apakah mas ${myBio.name} sudah menikah? ${myBio.isMarried} Kapan? ${myBio.when.age2})

// array of object

// let myFriend = [
//   {
//     name: "Nima",
//     age: 23,
//     address: "Tangerang",
//     isMarried: false
//   },
//   {
//     name: "Agung",
//     age: 32,
//     address: " ",
//     isMarried: true
//   }
// ]

// console.table(`Apakah mas dandi sudah menikah ? ${myBio.isMarried.age2}`)

// // array of object
// let myFriends = [
//     {
//         name: "Nima",
//         age: 23,
//         address: "Tangerang",
//         isMarried: true
//     },
//     {
//         name: "Agung",
//         age: 32,
//         address: "Jakarta",
//         isMarried: true,
//         istri: "angle"
//     },
//     {
//         name: "Krisna",
//         age: 15,
//         address: "Jakarta",
//         isMarried: true,
//         istri: "merri",
//         myCar: "supra"
//     }
// ]

// console.table(myFriends[0].myCar)

// console.table(myFriends[0].address)

// let example = document.getElementById("example")
// example.innerHTML += `
// <div>
//     <img src=${image} />
// </div>
// `


// 1. nilai dasar/awal
// 2. kondisi
// 3. initstatemen dan post statement

// for(let i = 1; i<=10; i++) {
//     console.table(`Hello world ${i}`)
// }

let blogs = []

function getData(event) {
  event.preventDefault()

  let title = document.getElementById('f-blog__title').value
  let content = document.getElementById('f-blog__content').value
  // let image = document.getElementById('f-blog__image').value
  /* Cara diatas akan membuat file menjadi tidak terbaca dan mengubahnya sebagai fake path,
  karena file input akan menghasilkan data berupa object(URL.createObjectURL()), sehingga perlu diubah seperti berikut:
  */
  let image = document.getElementById('f-blog__image').files

  image = URL.createObjectURL(image[0])

  console.table(image)

  let addBlog = {
    title,
    content,
    image,
    author: "Dandi Saputra",
    postedAt: new Date()
  }

  blogs.push(addBlog)

  console.log(blogs)
  showData()

}

function showData() {
  document.getElementById("contents").innerHTML = ""

  for (let i = 0; i <= blogs.length; i++) {
    document.getElementById("contents").innerHTML += `
        <div class="blog-list-item">
        <div class="blog-image">
          <img src="${blogs[i].image}" alt="" />
        </div>
        <div class="blog-content">
          <div class="btn-group">
            <button class="btn-edit">Edit Post</button>
            <button class="btn-post">Post Blog</button>
          </div>
          <h1>
            <a href="blog-detail.html" target="_blank"
              >${blogs[i].title}</a
            >
          </h1>
          <div class="detail-blog-content">
            ${blogs[i].postedAt} | ${blogs[i].author}
          </div>
          <p>${blogs[i].content}
          </p>
          <div style="float:right; margin: 10px">
            <p style="font-size: 15px; color:grey">1 minutes ago</p>
          </div>
        </div>
      </div>
        `
  }
}


// manipulation date time
// function getFullTime(time) {
//   // declaration variable
//   let years = time.getFullYear()
//   let monthIndex = time.getMonth()
//   let date = time.getDate()
//   let hour = time.getHours()
//   let minute = time.getMinutes()
  
//   const month = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]

//   return `${date} ${month[monthIndex]} ${years} ${hour}:${minute} WIB`
// }

// function getDistanceTime(time){
//     let timePost = time
//     let timeNow = new Date()

//     let milisecond = 1000 //milisecond

//     let distance = timeNow - timePost

//     let distanceDay = Math.floor(distance / (milisecond * 60 * 60 * 24)) //day 
//     let distanceHours = Math.floor(distance / (milisecond * 60 * 60)) //hours
//     let distanceMinute = Math.floor(distance / (milisecond * 60)) //minute
//     let distanceSecond = Math.floor(distance / milisecond) //second

//     if(distanceDay > 0) {
//       return `${distanceDay} Days Ago`
//     } else if(distanceHours > 0) {
//       return `${distanceHours} Hours Ago`
//     }else if(distanceMinute > 0) {
//       return `${distanceMinute} minutes Ago`
//     } else if(distanceSecond > 0) {
//       return `${distanceSecond} seconds Ago`
//     }
// }

// setInterval(() => {
//   showData()
// }, 1000)



// let posted = "2022-12-23 10:28:00"
// let generate = new Date("2022-12-22")

// console.log(typeof posted)
// console.log(typeof generate)
// // 1 tahun = 12
// // 1 bulan = 30
// // 1 hari = 24
// // 1 jam = 60
// // 1 menit = 60
// // 1 detik = 1000

// function getDistanceTime(time) {
//   let timePosting = time
//   let timeNow = new Date()

//   let distance = timeNow - timePosting

//   let daysDistance = Math.floor(distance / (24 * 60 * 60 * 1000))
//     if (daysDistance != 0) {
//         return daysDistance + ' Days Ago'
//     } else {
//         let hoursDistance = Math.floor(distance / (60 * 60 * 1000))
//         if (hoursDistance != 0) {
//             return hoursDistance + ' Hours Ago'
//         } else {
//             let minuteDistance = Math.floor(distance / (60 * 1000))
//             if (minuteDistance != 0) {
//                 return minuteDistance + ' Minutes Ago'
//             } else {
//                 let secondDistance = Math.floor(distance / 1000)
//                 if (secondDistance != 0)
//                     return secondDistance + ' sec'
//             }
//         }
//     }
// }