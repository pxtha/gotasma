import axios from 'axios'
import config from '../config'

class NetworkHelper {

    static requestPost(url, params, headers = null) {
        return NetworkHelper.requestHttp('POST', url, params, headers)
    }

    static requestGet(url, headers = null) {
        return NetworkHelper.requestHttp('GET', url, null, headers)
    }

    static requestPut(url, params, headers = null) {
        return NetworkHelper.requestHttp('PUT', url, params, headers)
    }

    static requestPatch(url, params, headers = null) {
        return NetworkHelper.requestHttp('PATCH', url, params, headers)
    }

    static requestDelete(url, params, headers = null) {
        return NetworkHelper.requestHttp('DELETE', url, params, headers)
    }

    static requestHttp(method, uri, params, headers) {
        return new Promise((resolve, reject) => {
            var options = {
                method,
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                }
            }

            if (params) {
                options.body = JSON.stringify(params)
            }
            if (headers) {
                options.headers['Authorization'] = 'Bearer ' + headers
            }
            const url = config.serverURI + uri
            fetch(url, options)
                .then((response) => {
                    response.json()
                        .then((body) => {
                            resolve(body)
                        })
                        .catch((error) => {
                            console.log(error)
                            reject('Can not connect to server')
                        })
                })
                .catch((error) => {
                    console.log(error)
                    reject('Can not connect to server')
                })
        })
    }

    static uploadFileWithProgress(url, file, key, params, onProgress, headers = null) {
        return new Promise((resolve, reject) => {
            var formData = new FormData()
            formData.append(key, file)
            for (var k in params) {
                formData.append(k, params[k])
            }
            axios.post(url, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': 'Bearer ' + headers
                    },
                    onUploadProgress: function(progressEvent) {
                        var percentCompleted = progressEvent.loaded / progressEvent.total
                        console.log(percentCompleted)
                        onProgress(percentCompleted)
                    }
                })
                .then((responseJson) => {
                    resolve({ statusCode: responseJson.status, body: responseJson.data })
                })
                .catch((error) => {
                    if (error.response !== undefined && error.response.status !== undefined) {
                        resolve({ statusCode: error.response.status, body: error.response.data })
                    } else {
                        reject('Can not connect to server')
                    }
                })
        })
    }

    static editFileWithProgress(url, file, key, params, onProgress, headers = null) {
        return new Promise((resolve, reject) => {
            var formData = new FormData()
            if (file) {
                formData.append(key, file)
            }
            for (var k in params) {
                formData.append(k, params[k])
            }
            axios.put(url, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': 'Bearer ' + headers
                    },
                    onUploadProgress: function(progressEvent) {
                        var percentCompleted = progressEvent.loaded / progressEvent.total
                        console.log(percentCompleted)
                        onProgress(percentCompleted)
                    }
                })
                .then((responseJson) => {
                    resolve({ statusCode: responseJson.status, body: responseJson.data })
                })
                .catch((error) => {
                    if (error.response !== undefined && error.response.status !== undefined) {
                        resolve({ statusCode: error.response.status, body: error.response.data })
                    } else {
                        reject('Can not connect to server')
                    }
                })
        })
    }
}

export default NetworkHelper