import React from "react";
import { QueryClient, QueryClientProvider } from "react-query";
import { Subjects } from "./pages/Subjects";

function App() {
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <Subjects />
    </QueryClientProvider>
  );
}

export default App;
