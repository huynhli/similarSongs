import { useEffect, useState } from 'react';
import '../App.css';
import { getAlbumRec, getArtistRec, getTrackRec } from './TanstackHelper';
import { useQuery } from '@tanstack/react-query';

type RecommendationWithArtist = {
  recName: string;
  artistName: string;
};

type RecommendationResponse = {
  type: 'track' | 'album' | 'artist';
  query: string;
  results: Record<string, RecommendationWithArtist[]> | RecommendationWithArtist[];
  error?: string;
};

type ResourceForQueryType = {
  type: 'track' | 'album' | 'artist';
  sanitizedLink: string;
};

function getRecs(resourceForQuery: ResourceForQueryType | null) {
  const queryFn =
    resourceForQuery?.type === 'album'
      ? () => getAlbumRec(resourceForQuery.sanitizedLink)
      : resourceForQuery?.type === 'artist'
      ? () => getArtistRec(resourceForQuery.sanitizedLink)
      : resourceForQuery?.type === 'track'
      ? () => getTrackRec(resourceForQuery.sanitizedLink)
      : () => Promise.resolve({ type: 'track', query: '', results: [] });

  return useQuery<RecommendationResponse>({
    queryKey: resourceForQuery ? [resourceForQuery.type, resourceForQuery.sanitizedLink] : [],
    queryFn: resourceForQuery ? queryFn : () => Promise.resolve({ type: 'track', query: '', results: [] }),
    enabled: !!resourceForQuery,
  });
}

export default function GeneratorPage() {
  const [resourceForQuery, setResourceForQuery] = useState<ResourceForQueryType | null>(null);
  const [link, setLink] = useState('');
  const [validationError, setValidationError] = useState<string | null>(null);

  const handleLinkChange = (e: React.ChangeEvent<HTMLInputElement>) => setLink(e.target.value);

  const sanitizeLink = (input: string): [string, string | null] => {
    const trimmed = input.trim();
    if (!trimmed) return [trimmed, 'No string entered'];
    if (!trimmed.startsWith('https://open.spotify.com/')) return [trimmed, 'Not a valid Spotify link: must start with "https://open.spotify.com/"'];
    return [trimmed, null];
  };

  const sanitizeValidateAndBackendCall = (inputLink: string) => {
    const [sanitizedLink, error] = sanitizeLink(inputLink);
    if (error) {
      setValidationError(error);
      return;
    }
    setValidationError(null);

    const split = sanitizedLink.split('/');
    const linkType = split[3]; // album / track / artist
    if (!['album', 'track', 'artist'].includes(linkType)) {
      setValidationError('Spotify link type not supported');
      return;
    }

    // Keep full link including ?si=...
    setResourceForQuery({ type: linkType as ResourceForQueryType['type'], sanitizedLink });
    setLink(sanitizedLink);
  };


  useEffect(() => {
    const url = new URLSearchParams(window.location.search);
    const linkQuery = url.get('link');
    if (linkQuery) {
      // decode once (URLSearchParams already decodes most things)
      sanitizeValidateAndBackendCall(linkQuery);
    }
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    sanitizeValidateAndBackendCall(link);
  };

  const { data, error, isLoading } = getRecs(resourceForQuery);

  const [tabIndexes, setTabIndexes] = useState<number[]>([0, 0, 0, 0]);

  const tags: string[] = data?.query.split(',').map(t => t.trim()) || [];
  const visibleTags = tags.slice(0, 5); // first 5 tags

  // Filter results for selected tag
  let recsForSelectedTag: RecommendationWithArtist[] = [];
  if (data?.results && !Array.isArray(data.results)) {
    const selectedTag = visibleTags[tabIndexes[3]]; // tabIndexes[3] for tag buttons
    recsForSelectedTag = data.results[selectedTag] || [];
  }

  return (
    // TODO make page not scrollable
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
        
        {/* results */}
        {resourceForQuery && (
          <div className='max-w-1000 w-[95%] grow min-mx-2 pb-2 border-b-1 border-white'>
            <div className='h-full flex flex-col w-full'>
              
              {/* Tab Headers */}
              <div className='h-20 flex justify-center'>
                <button className={`flex-1 text-center active:bg-[#222a3d] ${tabIndexes[0] === 0 ? "bg-[#222a3d] border-white  border-r-1" : "bg-[#090d14] hover:bg-black border-b-1 border-gray-600 border-b-white"} rounded-t-2xl flex items-center justify-center border-t-1 border-l-1`}
                  onClick={() => setTabIndexes(prev => [0, prev[1], prev[2], prev[3]])}
                >
                  LastFM
                </button>
                <button className={`flex-1 text-center active:bg-[#222a3d] ${tabIndexes[0] === 1 ? "bg-[#222a3d] border-white  border-l-1" : "bg-[#090d14] hover:bg-black border-b-1 border-gray-600 border-b-white"} rounded-t-2xl flex items-center justify-center border-t-1 border-r-1`}
                  onClick={() => setTabIndexes(prev => [1, prev[1], prev[2], prev[3]])}
                >
                  Explore
                </button>
              </div>

              {/* buttons */}
              <div className='w-full flex justify-around bg-[#222a3d] border-x-1 border-white'>
                {tabIndexes[0] === 0 ? 
                  <>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[1] === 0 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], 0, prev[2], prev[3]])}
                    >
                      Track Recs
                    </button>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[1] === 1 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], 1, prev[2], prev[3]])}
                    >
                      Similar Tracks
                    </button>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[1] === 2 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], 2, prev[2], prev[3]])}
                    >
                      Artist Recs
                    </button>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[1] === 3 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], 3, prev[2], prev[3]])}
                    >
                      Similar Artists
                    </button>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[1] === 4 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], 4, prev[2], prev[3]])}
                    >
                      Album Recs
                    </button>
                  </>
                : 
                  <>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[2] === 0 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], prev[1], 0, prev[3]])}
                    >
                      Acoustic Brainz Recs
                    </button>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[2] === 1 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], prev[1], 1, prev[3]])}
                    >
                      Discogs Recs
                    </button>
                    <button className={`mt-5 w-[15%] text-center active:bg-[#415178] ${tabIndexes[2] === 2 ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black "} ${window.innerWidth < 500 ? "text-sm p-1" : "p-2 text-lg"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], prev[1], 2, prev[3]])}
                    >
                      Deezer + Discogs Recs
                    </button>
                  </>
                }
              </div>
              
              {/* Spacer */}
              <div className='h-5 bg-[#222a3d] border-x-1 border-white'/>

              {/* Recs */}
              {/* Only show tag buttons for LastFM tab */}
              {tabIndexes[0] === 0 && visibleTags.length > 0 && (
                <div className='flex justify-center gap-2 pb-4 border-white border-x-1 bg-[#222a3d]'>
                  {visibleTags.map((tag, i) => (
                    <button
                      key={i}
                      className={`px-4 py-2 rounded-md text-white ${tabIndexes[3] === i ? "bg-[#415178]" : "bg-[#303c59] hover:bg-black"}`}
                      onClick={() => setTabIndexes(prev => [prev[0], prev[1], prev[2], i])}
                    >
                      {tag}
                    </button>
                  ))}
                </div>
              )}

              {/* LastFM recommendations grid */}
              {tabIndexes[0] === 0 && (tabIndexes[1] === 0 || tabIndexes[1] === 2 || tabIndexes[1] === 4) && recsForSelectedTag.length > 0 && (
                <div className='grid grid-cols-3 items-center justify-items-center gap-2 w-full bg-[#222a3d] border-x-1 border-white px-2'>
                  {recsForSelectedTag.map((rec, i) => (
                    <div key={i} className={`flex flex-col items-center shadow-md ${window.innerWidth < 500 ? "w-full" : "w-[80%]"} h-50`}>
                      <div className='flex flex-5 justify-center items-end rounded-t-xl border-white border-t-1 border-x-1 text-center w-full h-full bg-[#415178] pb-1'>{rec.recName}</div>
                      <div className='flex flex-4 justify-center items-start rounded-b-xl border-white border-b-1 border-x-1 text-center w-full h-full bg-[#3a486b] pt-1'>{rec.artistName}</div>
                    </div>
                  ))}
                </div>
              )}
              {tabIndexes[0] === 0 && (tabIndexes[1] === 1 || tabIndexes[1] === 3) && recsForSelectedTag.length > 0 && (
                <div className='grid grid-cols-3 items-center justify-items-center gap-2 w-full bg-[#222a3d] border-x-1 border-white px-2'>
                  {recsForSelectedTag.map((rec, i) => (
                    <div key={i} className={`flex flex-col items-center shadow-md ${window.innerWidth < 500 ? "w-full" : "w-[80%]"} h-50`}>
                      <div className='flex flex-5 justify-center items-end rounded-t-xl border-white border-t-1 border-x-1 text-center w-full h-full bg-[#415178] pb-1'>{rec.recName}</div>
                      <div className='flex flex-4 justify-center items-start rounded-b-xl border-white border-b-1 border-x-1 text-center w-full h-full bg-[#3a486b] pt-1'>{rec.artistName}</div>
                    </div>
                  ))}
                </div>
              )}

            </div>
          </div>
        )}
    </section>
  )
}