import '../App.css'

function DocsPage() {
    return (
        <div className="bg-blue-700 min-h-[calc(100vh-120px)] flex flex-row w-full justify-center items-center">
            <div className='bg-blue-700 w-20 h-20'></div>
            <h4 className='px-5 text-white text-center'>Currently this app only works for artists/songs/albums/playlists classified by 
                Spotify as hip hop, rap, pop, or rock. 
                <br />It was orignally designed to use more features of Spotify's Web API, but the most ideal features are now deprecated.
                <br />It is currently being updated, and more support will be updated in the future. 
                Sorry for the inconvenience.</h4>
                <div className='bg-blue-700 w-20 h-20'></div>
        </div>
    )
}

export default DocsPage