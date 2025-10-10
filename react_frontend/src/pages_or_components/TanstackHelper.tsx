import axiosHelper from "./AxiosHelper"

export const getTrackRec = async (id : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/track`+`/${id}`)
    return res.data
}

export const getTrackSimilarRec = async (id : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/track`+`/${id}`)
    return res.data
}

export const getArtistRec = async (id : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/artist`+`/${id}`)
    return res.data
}

export const getArtistSimilarRec = async (id : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/artist`+`/${id}`)
    return res.data
}

export const getAlbumRec = async (id : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/album`+`/${id}`)
    return res.data
}