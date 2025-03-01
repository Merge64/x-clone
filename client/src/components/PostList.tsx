import Post from './Post';

interface PostListProps {
  posts: any[];
  onRepost?: () => void;
  emptyMessage?: string;
}

function PostList({ posts, onRepost, emptyMessage = "No posts to display" }: PostListProps) {
  console.log("PostList received posts:", posts);
  
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