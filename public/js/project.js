
let dataProject = []

function getData(event) {
  event.preventDefault()

  let projectName = document.getElementById('project-name').value
  let projectStart = document.getElementById('project-start').value
  let projectEnd = document.getElementById('project-end').value
  let projectDesc = document.getElementById('project-description').value
  let projectTech = document.getElementsByName('project-tech')
  let projectImage = document.getElementById('project-image').files

  // ambil blob URL gambar
  projectImage = URL.createObjectURL(projectImage[0])

  let techChecked = []

  for (let i = 0; i < projectTech.length; i++) {
    if (projectTech[i].checked) {
      techChecked.push(projectTech[i].value)
    }
  }


  let addProject = {
    projectName,
    projectStart,
    projectEnd,
    projectDesc,
    projectImage,
    techChecked
  }

  dataProject.push(addProject)
  showData() // agar show data tidak dijalankan berulang

}

function showData() {
  document.getElementById("list-content").innerHTML = ""


  for (let i = 0; i < dataProject.length; i++) {
    listContent = document.getElementById('list-content')
    // currentContent = document.getelementById('post-' + i )

    listContent.innerHTML += `
    <div class="card card-post" id="post-${i}">
      <img src="${dataProject[i].projectImage}" alt="">
      <div class="">
        <div class="card-head">
          <h1 class="card-title__sm btn-link"><a href="detail-project.html">${dataProject[i].projectName}</a></h1>
          ${(function duration() {
        let string = ""
        if (dataProject[i].projectStart != "" && dataProject[i].projectEnd != "") {
          string = `<span class="card-subtitle__sm">${getDuration(dataProject[i].projectStart, dataProject[i].projectEnd)}</span>`
        }
        return string
      })()}
        </div>
        <div class="" style="height:100%;">
          <p class="card-desc__sm">${dataProject[i].projectDesc}</p>
          <ul class="list-items-sm">
            ${(function icon() {
        let string = ""
        for (let j = 0; j < dataProject[i].techChecked.length; j++) {
          string += `<li><img src="assets/img/icon/logo-${dataProject[i].techChecked[j]}.svg" alt="Item Icon"></li>`
        }
        return string
      })()}
            </ul>
        </div>
        <div class="btn-group" style="margin-top:24px;">
          <a class="btn btn-primary btn-sm btn-full" href="single-project.html">Edit</a>
          <span onclick="document.getElementById('post-${i}').remove();" class="btn btn-danger btn-sm btn-full" >Delete</span>
        </div>
      </div>
    </div>
    `
  }
}


function getDuration(start, end) {
  let projectStart = new Date(start)
  let projectEnd = new Date(end)
  let range = projectEnd - projectStart
  let monthRange = Math.floor(range / (30 * 24 * 60 * 60 * 1000))
  if (monthRange < 0) {
    return ''
  }
  if (monthRange > 0) {
    return 'Duration : ' + monthRange + ' Month'
  } else {
    let weekRange = Math.floor(range / (7 * 24 * 60 * 60 * 1000))
    if (weekRange != 0) {
      return 'Duration : ' + weekRange + ' Week Left'
    } else {
      let weekRange = Math.floor(range / (7 * 24 * 60 * 60 * 1000))
      if (weekRange != 0) {
        return 'Duration : ' + weekRange + ' Week Left'
      } else {
        let daysRange = Math.floor(range / (24 * 60 * 60 * 1000))
        if (daysRange != 0) {
          return 'Duration : ' + daysRange + ' Days Left'
        } else {
          return 'Duration : Today'
        }
      }
    }
  }
}
