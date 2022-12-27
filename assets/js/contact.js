function submitData() {
  let name = document.getElementById('f-name').value
  let email = document.getElementById('f-email').value
  let phone = document.getElementById('f-phone').value
  let subject = document.getElementById('f-subject').value
  let message = document.getElementById('f-message').value

  if (name == '') {
    return alert('name Form is required')
  }
  else if (email == '') {
    return alert('Email Form is required')
  }
  else if (phone == '') {
    return alert('Phone Number Form is required')
  }
  else if (subject == '') {
    return alert('Subject Form is required')
  }
  else if (message == '') {
    return alert('Message Form is required')
  }

  // For sending Mail
  let emailRecipient = "wahyudi.chrisdianto@gmail.com"

  let mailTo = document.createElement('a')
  mailTo.href = `mailto:${emailRecipient}?subject=${subject}&body=Halo Perkenalkan nama saya ${name}%0D%0A%0D%0A${message},%0D%0Asilahkan hubungi saya di ${phone}%0D%0A%0D%0ATerima kasih.`
  mailTo.click()
}