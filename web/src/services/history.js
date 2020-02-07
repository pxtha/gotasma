import NetworkHelper from './NetworkHelper'

export const addHistory = async(snapshot) => {
    return await NetworkHelper.requestPost('/history', snapshot)
}

export const getHistoryByID = async(id) => {
    return await NetworkHelper.requestGet('/history?projectId=' + id)
}