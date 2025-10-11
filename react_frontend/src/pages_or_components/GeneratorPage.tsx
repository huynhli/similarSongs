import { useEffect, useState } from 'react'
import '../App.css'
import { getAlbumRec, getArtistRec, getTrackRec } from './TanstackHelper';
import { useQuery } from '@tanstack/react-query'

type RecommendationWithArtist = {
  recName: string;
  artistName: string;
};

type RecommendationResponse = {
  type: "track" | "album" | "artist";
  query: string;
  results:
    | Record<string, RecommendationWithArtist[]>
    | RecommendationWithArtist[];
  error?: string;
};

type resourceForQueryType = {
  type: string,
  sanitizedLink: string,
}

function getRecs(resourceForQuery: resourceForQueryType | null) {
  const queryFnHelp = resourceForQuery?.type === "album"
    ? () => getAlbumRec(resourceForQuery.sanitizedLink)
    : resourceForQuery?.type === "artist"
      ? () => getArtistRec(resourceForQuery.sanitizedLink)
      : resourceForQuery?.type === "track"
        ? () => getTrackRec(resourceForQuery.sanitizedLink)
        : () => Promise.resolve({ type: "track", query: "", results: [] }); 

  return useQuery<RecommendationResponse>({
    queryKey: resourceForQuery ? [resourceForQuery.type, resourceForQuery.sanitizedLink] : [],
    queryFn: resourceForQuery ? queryFnHelp! : () => Promise.resolve([]),
    enabled: !!resourceForQuery // --> only runs if theres a resourceForQuery, also enables auto refetching
  })
}

export default function GeneratorPage() {
  const [resourceForQuery, setResourceForQuery] = useState<resourceForQueryType | null>(null)
  const [link, setLink] = useState("");

  const handleLinkChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setLink(event.target.value);
  }

  useEffect(() => {
    // simple + stateless
    const url = new URLSearchParams(window.location.search)
    const linkQuery = url.get('link')

    if (linkQuery){
      const decodedLink = decodeURIComponent(linkQuery)
      setLink(decodedLink) // change box text, and state
      sanitizeValidateAndBackendCall(decodedLink) // on page load, use decoded link and not state, bc state updates async
    }
  }, [])

  const handleSubmit = (e: React.FormEvent) => {
    e?.preventDefault()
    sanitizeValidateAndBackendCall(link)
  }
  
  const sanitizeLink = (link : string) : [string, string | null] => {
    if (!link) return [link, "No string entered"]
    const trimmedLowerLink = link.trim().toLowerCase()
    if (!trimmedLowerLink.startsWith('https://open.spotify.com/')) {
      return [trimmedLowerLink, 'Not a valid spotify link: Does not start with "https://open.spotify.com/"']
    }
    return [trimmedLowerLink, null]
  }

  const [validationError, setValidationError] = useState<string | null>(null)

  const sanitizeValidateAndBackendCall = (decodedLink : string) => {
    const [sanitizedLink, error] = sanitizeLink(decodedLink)
    if (error) {
      setValidationError(error)
      return
    }
    setValidationError(null)

    console.log('Link passed valididation: ' + sanitizedLink)
    console.log('Requesting info...')
    const splitSanitizedLink = sanitizedLink?.split('/')
    const id = splitSanitizedLink[4] 
    const linkType = splitSanitizedLink[3]
    if (!["album", "track", "artist"].includes(linkType)){
      return console.log('Spotify link type not supported. Please use a different link.')
    }
    
    // runs getRecs implicitly
    setResourceForQuery({type: linkType, sanitizedLink: sanitizedLink})
  }

  const { data, error, isLoading } = getRecs(resourceForQuery)

  const [ tabIndexes, setTabIndexes ] = useState<number[]>([0, 0, 0])

  // Flatten results for easier mapping
  let recs: RecommendationWithArtist[] = []
  if (data?.results) {
    if (Array.isArray(data.results)) {
      recs = data.results
    } else {
      // results is Record<string, RecommendationWithArtist[]>
      recs = Object.values(data.results).flat()
    }
  }

  return (
    <section className='min-h-[100vh] mb-2 w-full flex flex-col justify-center items-center'>
        
        {/* Input */}
        <div className={`${!data?.results ? "h-full" : "flex-1"} min-h-80 max-h-120 w-full flex justify-center items-center flex-col`}>
          <form className={`max-w-[1600px] mx-[10%] text-black flex ${window.innerWidth < 600 ? "flex-col" : "flex-row"}`} onSubmit={handleSubmit}>
            <div className='flex justify-center items-center'>
              <input
                className='flex px-2 mx-2 w-auto max-w-600 min-w-100 bg-white'
                type='text'
                value={link}
                onChange={handleLinkChange}
                placeholder='Enter a spotify link here'
              />
            </div>
            <div className='flex justify-center items-center'>
              <button className='px-2 mx-2 rounded-md w-30 border-black border-1 bg-gray-300 hover:bg-gray-200 active:bg-gray-100' type="submit">Generate</button>
            </div>
          </form>
          {validationError && (
            <div className='my-2 flex'>{validationError}</div>
          )}
        </div>

        {/* Results */}
        {resourceForQuery && (
          <div className='max-w-1000 w-[95%] grow min-mx-2 pb-2 border-b-1 border-white'>
            
            {/* Loading/Error/No Results */}
            {isLoading && (
              <div className='h-full flex justify-center items-center text-white'>Loading recommendations...</div>
            )}
            {error && (
              <div className='h-full flex justify-center items-center text-red-500'>
                {error instanceof Error ? error.message : "An error occurred while fetching recommendations"}
              </div>
            )}
            {!isLoading && !error && recs.length === 0 && (
              <div className='h-full flex justify-center items-center text-white'>0 recommendations found</div>
            )}

            {/* Only render tabs and recs if data is loaded */}
            {!isLoading && !error && recs.length > 0 && (
              <div className='h-full flex flex-col w-full'>
                {/* Tab Headers */}
                <div className='h-20 flex justify-center'>
                  <button className={`flex-1 text-center hover:bg-black active:bg-[#222a3d] ${tabIndexes[0] === 0 ? "bg-[#222a3d] border-white  border-r-1" : "bg-[#090d14] border-b-1 border-gray-600 border-b-white"} rounded-t-2xl flex items-center justify-center border-t-1 border-l-1`}
                    onClick={() => setTabIndexes(prev => [0, prev[1], prev[2]])}
                  >
                    LastFM
                  </button>
                  <button className={`flex-1 text-center hover:bg-black active:bg-[#222a3d] ${tabIndexes[0] === 1 ? "bg-[#222a3d] border-white  border-l-1" : "bg-[#090d14] border-b-1 border-gray-600 border-b-white"} rounded-t-2xl flex items-center justify-center border-t-1 border-r-1`}
                    onClick={() => setTabIndexes(prev => [1, prev[1], prev[2]])}
                  >
                    Explore
                  </button>
                </div>

                {/* Buttons */}
                <div className='w-full flex justify-around bg-[#222a3d] border-x-1 border-white'>
                  {/* ... keep your button JSX here ... */}
                </div>

                {/* Spacer */}
                <div className='h-5 bg-[#222a3d] border-x-1 border-white'/>

                {/* Recommendation Grid */}
                <div className='grid grid-cols-3 items-center justify-items-center gap-2 w-full bg-[#222a3d] border-x-1 border-white px-2'>
                  {recs.map((rec, i) => (
                    <div key={i} className={`flex flex-col items-center shadow-md ${window.innerWidth < 500 ? "w-full" : "w-[80%]"} ${tabIndexes[1] === 2 || tabIndexes[1] === 3 ? "h-30" : "h-50"}`}>
                      {/* artist cases */}
                      {tabIndexes[1] === 2 || tabIndexes[1] === 3 ? 
                        <div className='flex flex-4 justify-center items-center rounded-xl border-white border-1 text-center w-full h-full bg-[#3a486b] pt-1'>
                          {rec.artistName}
                        </div>
                      :
                        <>
                          <div className='flex flex-5 justify-center items-end rounded-t-xl border-white border-t-1 border-x-1 text-center w-full h-full bg-[#415178] pb-1'>{rec.recName}</div>
                          <div className='flex flex-4 justify-center items-start rounded-b-xl border-white border-b-1 border-x-1 text-center w-full h-full bg-[#3a486b] pt-1'>{rec.artistName}</div>
                        </>
                      }
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
    </section>
  )
}
