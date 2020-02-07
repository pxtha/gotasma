const urlParser = document.createElement('a')
import moment from 'moment'

export function domain(url) {
    urlParser.href = url
    return urlParser.hostname
}

export function count(arr) {
    return arr.length
}

export function changeDate(value) {
    return (value / 86400000)
}

export function prettyDate(date) {
    var a = new Date(date)
    return a.toDateString()
}

export function momentNormalDate(date) {
    return moment(date).format('MMMM Do YYYY')
}

export function momentDetailDate(date) {
    return moment(date).format('MMM Do YYYY, h:mm:ss a')
}

export function pluralize(time, label) {
    if (time === 1) {
        return time + label
    }

    return time + label + 's'
}