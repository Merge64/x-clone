import Post from './Post';

interface PostListProps {
  posts: any[];
  onRepost?: () => void;
  emptyMessage?: string;
  isLoading: boolean; // Add isLoading to the props
}

function PostList({ posts, onRepost, emptyMessage = "No posts to display", isLoading }: PostListProps) {
  console.log("PostList received posts:", posts);

  // Display loading state if isLoading is true
  if (isLoading) {
    return (
      <div className="p-6 text-center text-gray-500">
        <p>Loading...</p>
      </div>
    );
  }

  // Check if posts exists and is an array
  if (!posts || !Array.isArray(posts) || posts.length === 0) {
    return (
      <div className="p-6 text-center text-gray-500">
        <p>{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div>
      {posts.map((post) => (
        <Post key={post.id} post={post} onRepost={onRepost} />
      ))}
    </div>
  );
}

export default PostList;
