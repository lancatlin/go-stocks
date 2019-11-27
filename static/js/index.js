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
