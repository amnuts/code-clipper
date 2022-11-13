<script>
    import { onMount } from 'svelte';
    import TagList from "./TagList.svelte";
    import SnippetsList from "./SnippetsList.svelte";
    import SnippetView from "./SnippetView.svelte";
    import { GetSnippets, GetTags } from "../../wailsjs/go/main/App.js";

    let snippets, tags, totalSnippets, selectedTag, selectedSnippet, selectedSnippetId;

    const fetchSnippets = async () => {
        snippets = await GetSnippets();
        tags = await GetTags();
        totalSnippets = snippets ? Object.keys(snippets).length : 0;
    }

    onMount(fetchSnippets);

    function handleTagSelection(event) {
        selectedTag = event.detail.tag;
    }

    function handleSnippetSelection(event) {
        selectedSnippet = snippets[event.detail.snippet];
        selectedSnippetId = event.detail.snippet;
    }

    function handleSnippetSave(event) {
        (fetchSnippets)();
        if (event.detail.snippet === null) {
            selectedSnippetId = null;
        }
    }
</script>

<main>
    <TagList allTags={tags} totalSnippets={totalSnippets} on:message={handleTagSelection} />
    <SnippetsList {...{snippets, totalSnippets, selectedSnippetId, selectedTag}} on:message={handleSnippetSelection} />
    <SnippetView selectedSnippet={selectedSnippet} on:message={handleSnippetSave} />
</main>

<style>
    main {
        width: 100vw;
        height: 100vh;
        padding: 0;
        margin: 0;
        border: 0;
        background-color: #333333;
        font-size: 80%;
        display: flex;
        flex-direction: row;
    }
</style>