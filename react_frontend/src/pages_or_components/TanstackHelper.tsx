import axiosHelper from "./AxiosHelper"

export const getTrackRec = async (link : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/track`, { params: { link: link }})
    return res.data
}

export const getTrackSimilarRec = async (link : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/track`, { params: { link: link }})
    return res.data
}

export const getArtistRec = async (link : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/artist`, { params: { link: link }})
    return res.data
}

export const getArtistSimilarRec = async (link : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/artist`, { params: { link: link }})
    return res.data
}

export const getAlbumRec = async (link : string) => {
    const res = await axiosHelper.get(`/lastfm`+`/album`, { params: { link: link }})
    return res.data
}