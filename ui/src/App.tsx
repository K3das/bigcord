import type {Component} from 'solid-js';
import {APIClient} from "./utils/api"
import ChannelStates from "./components/ChannelStates";
import GrafanaPanel from "./components/GrafanaPanel";
import AddChannels from "./components/AddChannels";

const App: Component = () => {
    const api = new APIClient(import.meta.env.VITE_API_BASE_URL ?? "/api")
    
    return (
        <div class={"max-w-7xl mx-auto px-5 py-5 md:py-10"}>
            <AddChannels api={api}></AddChannels>
            <ChannelStates api={api}></ChannelStates>
            <GrafanaPanel dashboard={"bigcord-provisioned-dashboard-main"} refresh={"5s"} panelId={57}/>
            <span
                class={"my-3 block text-sm text-gray-700 sm:text-center"}>Bigcord @{import.meta.env.VITE_COMMIT_HASH ?? "DEV_BUILD"}</span>
        </div>
    );
};

export default App;
