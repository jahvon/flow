import { useState } from "react";
import reactLogo from "./assets/react.svg";
import { invoke } from "@tauri-apps/api/core";
import "./App.css";

function App() {
  const [greetMsg, setGreetMsg] = useState("");
  const [name, setName] = useState("");
  const [workspaces, setWorkspaces] = useState([])



  async function greet() {
    // Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
    // setGreetMsg(await invoke("greet", { name }));
      console.log(("hello"))
      await invoke("workspaces")
          .then((message) => setWorkspaces(message.workspaces))
          .catch((error) => console.error(error));
  }


  return (
    <main className="container">
      <h1>Welcome to Tauri + React</h1>
        <button onClick={async () => {
            console.log(("hello"))
            await invoke("workspaces").then((message) => setWorkspaces(message.workspaces))
                .catch((error) => console.error(error));
        }
        }>Workspaces</button>
        {workspaces.map((workspace, i) => {
                return (
                    <div key={i}>
                        <p>{workspace.displayName}</p>
                        {workspace.description && <p>{workspace.description}</p>}
                    </div>

                )
        })
        }

        <form
            className="row"
        onSubmit={(e) => {
          e.preventDefault();
          greet();
        }}
      >
        <input
          id="greet-input"
          onChange={(e) => setName(e.currentTarget.value)}
          placeholder="Enter a name..."
        />
        <button type="submit">Greet</button>
      </form>
      <p>{greetMsg}</p>
    </main>
  );
}

export default App;
