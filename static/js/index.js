function save(id) {
    let cookie = Cookies.get('id')
    if (cookie == undefined) {
        Cookies.set('id', id, { expires: 365 })
    } else {
        Cookies.set('id', `${id}&${cookie}`, { expires: 365 })
    }
    window.location.href = '/'
}

function remove(id) {
    let stocks = Cookies.get('id').split('&')
    for (i in stocks) {
        if (stocks[i] == id) {
            stocks.splice(i, 1)
            break
        }
    }
    Cookies.set('id', stocks.join('&'), { expires: 365 })
    document.getElementById(`tr-${id}`).remove()
}

function search(input) {
    let q = input.value
    console.log(q)
    let xhttp = new XMLHttpRequest()
    xhttp.open("GET", `/api/search?q=${q}`)
    xhttp.onreadystatechange = function () {
        if (this.readyState != 4) {
            return
        }
        if (this.status == 200) {
            console.log(this.status)
            update(parse(this.responseText))
        }
    }
    xhttp.send()
}

function parse(text) {
    return JSON.parse(text)
}

function update(stocks) {
    let datalist = document.getElementById("data")
    while (datalist.firstChild) {
        datalist.removeChild(datalist.firstChild);
    }
    for (stock of stocks) {
        let elem = document.createElement("option")
        elem.value = stock.ID
        elem.innerText = stock.ID + stock.Name
        datalist.appendChild(elem)
    }
}
