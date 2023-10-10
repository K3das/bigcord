import {Channel, ChannelType, JobState} from "../utils/api";
import {createEffect, createMemo, createSignal, For} from "solid-js";

// TODO: Fix channel updates
const ChannelList = (props: { channels: Channel[], onUpdate: (channels: Set<string>) => void }) => {
    const [selectedChannels, setSelectedChannels] = createSignal<Set<string>>(new Set());
    const [selectedCategories, setSelectedCategories] = createSignal<Set<string>>(new Set());

    createEffect(() => {
        props.onUpdate(selectedChannels())
    })
    
    const channels = createMemo(() => {
        setSelectedChannels(new Set<string>());
        setSelectedCategories(new Set<string>());
        return props.channels.filter((channel) =>
            channel.job_state != JobState.RUNNING &&
            channel.has_access)
    })
    
    const toggleChannelSelection = (channelId: string) => {
        setSelectedChannels(prevSelected => {
            const newSelected = new Set(prevSelected);
            if (newSelected.has(channelId)) {
                newSelected.delete(channelId);
            } else {
                newSelected.add(channelId);
            }
            return newSelected;
        });
    };

    const toggleCategorySelection = (categoryId: string) => {
        let value = true;
        setSelectedCategories(prevSelected => {
            const newSelected = new Set(prevSelected);
            if (newSelected.has(categoryId)) {
                newSelected.delete(categoryId);
                value = false;
            } else {
                newSelected.add(categoryId);
            }
            return newSelected;
        });

        setSelectedChannels(prevSelected => {
            const newSelected = new Set(prevSelected);
            channels()
                .filter(channel => channel.parent_id === categoryId)
                .forEach((channel) => {
                    if (value) {
                        newSelected.add(channel.id);
                    } else {
                        newSelected.delete(channel.id);
                    }
                })
            return newSelected;
        });
    };


    const renderChannels = (parent?: string) => {
        return (
            <For each={channels().filter(channel => channel.parent_id === parent)}>{(channel, _i) => (
                <div title={ChannelType[channel.type]} class={"font-mono"}>
                    {channel.type == ChannelType.GUILD_CATEGORY ? (
                        <><label class="flex items-center space-x-2">
                            <input
                                type="checkbox"
                                checked={selectedCategories().has(channel.id)}
                                onChange={() => toggleCategorySelection(channel.id)}/>
                            <span>{channel.name}:</span>
                        </label>
                            <div class={"ml-4"}>
                                {renderChannels(channel.id)}
                            </div>
                        </>
                    ) : (
                        <label class="flex items-center space-x-2">
                            <input
                                type="checkbox"
                                checked={selectedChannels().has(channel.id)}
                                onChange={() => toggleChannelSelection(channel.id)}
                                disabled={selectedCategories().has(channel.parent_id)}
                            />
                            <span class={"p-2"}>#{channel.name}</span>
                        </label>
                    )}
                </div>
            )}</For>
        )
    };

    return renderChannels();
};

export default ChannelList;
