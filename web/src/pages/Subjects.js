import React from 'react';
import {client} from '../utils/client'
import {useQuery} from 'react-query';


export function Subjects() {
    const query = useQuery("subjects", () => client.get("/subjects"))

    if (query.isLoading) {
        return <div>Loading...</div>
    }

    if (query.isError) {
        return <div>Error {query.error.message}</div>
    }

    return (

        <div>
            {query.data.map(s => <div>{s.Name}</div>)}
        </div>
    );
}

