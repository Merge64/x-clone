import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, MoreVertical, MessageCircle, Heart, Share, Repeat } from 'lucide-react';
import { 
  getPostById, 
  getUserInfo, 
  repost, 
  toggleLike, 
  deletePost, 
  editPost, 
  editQuote, 
  checkIfLiked, 
  checkIfReposted, 
  addComment, 
  getComments,
  getRepostsCount,
  getCommentsCount,
  getLikesCount
} from '../utils/api';
import { PostData } from '../types/post';
import PostOptions from '../components/posts/PostOptions';
import CommentModal from '../components/posts/CommentModal';
import QuoteModal from '../components/posts/QuoteModal';
import PostContent from '../components/posts/PostContent';
import Navbar from './navbar/Navbar';
import Post from '../components/posts/Post';

const defaultPost: PostData = {
    id: 0,
    created_at: new Date().toISOString(),
    userid: 0,
    username: '',
    nickname: '',
    body: '',
    likes_count: 0,
    reposts_count: 0,
    comments_count: 0,
    is_repost: false,
    parent_id: undefined,
    quote: '',
    parent_post: undefined,

};

function PostDetailPage() {
    const { username, postId } = useParams<{ username: string; postId: string }>();
    const navigate = useNavigate();
    const [post, setPost] = useState<PostData>(defaultPost);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isLiked, setIsLiked] = useState(false);
    const [isReposted, setIsReposted] = useState(false);
    const [showRepostOptions, setShowRepostOptions] = useState(false);
    const [isReposting, setIsReposting] = useState(false);
    const [isUndoingRepost, setIsUndoingRepost] = useState(false);
    const [currentUser, setCurrentUser] = useState<any>(null);
    const [showOptions, setShowOptions] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [editContent, setEditContent] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [editError, setEditError] = useState('');
    const [showQuoteModal, setShowQuoteModal] = useState(false);
    const [showCommentModal, setShowCommentModal] = useState(false);
    const [newComment, setNewComment] = useState('');
    const [isAddingComment, setIsAddingComment] = useState(false);
    const [comments, setComments] = useState<any[]>([]);
    const [loadingComments, setLoadingComments] = useState(true);
    const [likesCount, setLikesCount] = useState(0);
    const [repostsCount, setRepostsCount] = useState(0);
    const [commentsCount, setCommentsCount] = useState(0);
    const [isLikeProcessing, setIsLikeProcessing] = useState(false);

    const MAX_CHARS = 280;
    const isQuoteRepost = Boolean(post.is_repost && post.parent_id !== null && post.quote);
    const isSimpleRepost = Boolean(post.is_repost && post.parent_id !== null && !post.quote);
    const isCurrentUserAuthor = post.username === currentUser?.username;

    useEffect(() => {
        const fetchLikesCount = async () => {
            if (post.is_repost && (post.quote === null || post.quote === undefined || post.quote === "")) {
                try {
                    const count = await getLikesCount(Number(post.parent_post?.id));
                    setLikesCount(count);
                } catch (error) {
                    console.error('Error fetching likes count:', error);
                }
            } else {
                try {
                    const count = await getLikesCount(Number(post.id));
                    setLikesCount(count);
                    console.log(post.parent_post)

                } catch (error) {
                    console.error('Error fetching likes count:', error);
                }
            }
        };

        fetchLikesCount();
    }, [post.is_repost, post.parent_id, post.id, post.parent_post?.id]);

    useEffect(() => {
        const fetchRepostsCount = async () => {
            if (post.is_repost && (post.quote === null || post.quote === undefined || post.quote === "")) {
                try {
                    const count = await getRepostsCount(Number(post.parent_post?.id));
                    setRepostsCount(count);
                } catch (error) {
                    console.error('Error fetching reposts count:', error);
                }
            } else {
                try {
                    const count = await getRepostsCount(Number(post.id));
                    setRepostsCount(count);
                } catch (error) {
                    console.error('Error fetching reposts count:', error);
                }
            }
        };

        fetchRepostsCount();
    }, [post.is_repost, post.parent_id, post.id, post.parent_post?.id]);

    useEffect(() => {
        const fetchCommentsCount = async () => {
          // Simplified condition
          if (post.is_repost && !post.quote) {
            try {
              const count = await getCommentsCount(Number(post.parent_id));
              setCommentsCount(count);
            } catch (error) {
              console.error('Error fetching comments count:', error);
            }
          } else {
            try {
              const count = await getCommentsCount(Number(post.id));
              setCommentsCount(count);
            } catch (error) {
              console.error('Error fetching comments count:', error);
            }
          }
        };
      
        fetchCommentsCount();
      }, [post.is_repost, post.parent_id, post.id]); // parent_id instead of parent_post.id

    useEffect(() => {
        const fetchComments = async () => {
          if (!post.id) return;
          
          setLoadingComments(true);
          try {
            // Use parent_id for simple reposts
            if (post.is_repost && !post.quote) {
              const fetchedComments = await getComments(Number(post.parent_id));
              setComments(fetchedComments);
            } else {
              const fetchedComments = await getComments(Number(post.id));
              setComments(fetchedComments);
            }
          } catch (error) {
            console.error('Error fetching comments:', error);
          } finally {
            setLoadingComments(false);
          }
        };
      
        fetchComments();
      }, [post.id, post.is_repost, post.parent_id]); // Use parent_id in dependency array

    useEffect(() => {
        const checkInteractions = async () => {
            if (post.is_repost && (post.quote === null || post.quote === undefined || post.quote === "")) {
                try {
                    const [likedStatus, repostedStatus] = await Promise.all([
                        checkIfLiked(Number(post.parent_id)),
                        checkIfReposted(Number(post.parent_id))
                    ]);
                    setIsLiked(likedStatus);
                    setIsReposted(repostedStatus);
                } catch (error) {
                    console.error('Error checking post interactions:', error);
                }
            } else {
                try {
                    const [likedStatus, repostedStatus] = await Promise.all([
                        checkIfLiked(Number(post.id)),
                        checkIfReposted(Number(post.id))
                    ]);
                    setIsLiked(likedStatus);
                    setIsReposted(repostedStatus);
                } catch (error) {
                    console.error('Error checking post interactions:', error);
                }
            }
        };

        if (currentUser && post.id) {
            checkInteractions();
        }
    }, [post.id, post.is_repost, post.parent_id, currentUser, post.parent_post?.id]);

    useEffect(() => {
        const handleClickOutside = () => {
            if (showRepostOptions) {
                setShowRepostOptions(false);
            }
        };

        document.addEventListener('click', handleClickOutside);
        return () => {
            document.removeEventListener('click', handleClickOutside);
        };
    }, [showRepostOptions]);

    const handleCommentAdded = async () => { 
        if (post.id) {
          // Always fetch from parent_id for simple reposts
          if (post.is_repost && !post.quote) {
            const fetchedComments = await getComments(Number(post.parent_id));
            setComments(fetchedComments);
          } else {
            const fetchedComments = await getComments(Number(post.id));
            setComments(fetchedComments);
          }
        }
      };

    const formatDate = (isoString: string | number | Date) => {
        const date = new Date(isoString);
        return date.toLocaleTimeString('en-US', {
            hour: 'numeric',
            minute: '2-digit',
            hour12: true,
        }) + ' · ' + date.toLocaleString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        });
    };

    useEffect(() => {
        const fetchPost = async () => {
            if (!username || !postId) {
                setError('Invalid URL parameters');
                setLoading(false);
                return;
            }

            try {
                const postData = await getPostById(postId);
                if (!postData) {
                    setError('Post not found');
                    return;
                }
                setPost(postData);
                setEditContent(isQuoteRepost ? postData.quote || '' : postData.body);
            } catch (err) {
                setError('Failed to load post');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchPost();
    }, [username, postId]);

    useEffect(() => {
        const fetchUserInfo = async () => {
            try {
                const info = await getUserInfo();
                setCurrentUser(info);
            } catch (error) {
                console.error('Error fetching user info:', error);
            }
        };

        fetchUserInfo();
    }, []);

    const toggleOptions = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setShowOptions(!showOptions);
    };

    const handleEdit = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setIsEditing(true);
        setShowOptions(false);
    };

    const handleCancelEdit = () => {
        setIsEditing(false);
        setEditError('');
    };

    const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const newContent = e.target.value;
        if (newContent.length <= MAX_CHARS) {
            setEditContent(newContent);
        }
    };

    const handleCommentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const text = e.target.value;
        if (text.length <= MAX_CHARS) {
            setNewComment(text);
        }
    };

    const submitComment = async () => {
        if (!newComment.trim() || isAddingComment) return;
      
        setIsAddingComment(true);
        try {
          // Submit to parent post if it's a simple repost
          const targetPostId = post.is_repost && !post.quote 
            ? Number(post.parent_id) 
            : Number(post.id);
          
          await addComment(targetPostId, newComment);
          setCommentsCount(prev => prev + 1);
          setNewComment('');
          await handleCommentAdded();
        } catch (error) {
          console.error('Error adding comment:', error);
          alert('Failed to post comment. Please try again.');
        } finally {
          setIsAddingComment(false);
        }
      };

    const handleSaveEdit = async () => {
        if (!editContent.trim()) {
            setEditError('Post cannot be empty');
            return;
        }

        setIsSubmitting(true);
        setEditError('');

        try {
            if (isQuoteRepost) {
                await editQuote(post.id.toString(), editContent);
            } else {
                await editPost(post.id.toString(), editContent);
            }

            setPost(prev => ({
                ...prev,
                [isQuoteRepost ? 'quote' : 'body']: editContent
            }));
            setIsEditing(false);
        } catch (error) {
            console.error(`Error editing ${isQuoteRepost ? 'quote' : 'post'}:`, error);
            setEditError(`Failed to edit ${isQuoteRepost ? 'quote' : 'post'}. Please try again.`);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleDelete = async (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();

        if (isDeleting) return;

        if (window.confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
            setIsDeleting(true);
            try {
                await deletePost(post.id.toString());
                navigate(-1);
            } catch (error) {
                console.error('Error deleting post:', error);
                alert('Failed to delete post. Please try again.');
            } finally {
                setIsDeleting(false);
                setShowOptions(false);
            }
        } else {
            setShowOptions(false);
        }
    };

    const handleLikeToggle = async (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();

        if (isLikeProcessing) return;

        setIsLikeProcessing(true);
        const previousLikeState = isLiked;
        const previousCount = likesCount;

        try {
            setIsLiked(!previousLikeState);
            setLikesCount(prevCount => previousLikeState ? prevCount - 1 : prevCount + 1);

            if (post.is_repost && (post.quote === null || post.quote === undefined || post.quote === "")) {
                await toggleLike(Number(post.parent_id));
                
                const currentLikeStatus = await checkIfLiked(Number(post.parent_id));
                const currentLikesCount = await getLikesCount(Number(post.parent_id));
                setIsLiked(currentLikeStatus);
                setLikesCount(currentLikesCount);
            } else {
                await toggleLike(Number(post.id));
                
                const currentLikeStatus = await checkIfLiked(Number(post.id));
                const currentLikesCount = await getLikesCount(Number(post.id));
                setIsLiked(currentLikeStatus);
                setLikesCount(currentLikesCount);
            }
        } catch (error) {
            console.error('Error toggling like:', error);
            setIsLiked(previousLikeState);
            setLikesCount(previousCount);
        } finally {
            setIsLikeProcessing(false);
        }
    };

    const toggleRepostOptions = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (isReposted) {
            handleUndoRepost(e);
        } else {
            setShowRepostOptions(!showRepostOptions);
        }
    };

    const handleSimpleRepost = async (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();

        if (isReposting) return;

        setIsReposting(true);
        const prevRepostedState = isReposted;
        const prevRepostsCount = repostsCount;

        try {
            setIsReposted(true);
            setRepostsCount(prev => prev + 1);

            if (post.is_repost) {
                await repost(Number(post.parent_id), "");
            } else {
                await repost(Number(post.id), "");
            }
            setShowRepostOptions(false);
        } catch (error) {
            console.error('Error reposting:', error);
            setIsReposted(prevRepostedState);
            setRepostsCount(prevRepostsCount);
            alert('Failed to repost. Please try again.');
        } finally {
            setIsReposting(false);
        }
    };

    const handleUndoRepost = async (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();

        if (isUndoingRepost) return;

        setIsUndoingRepost(true);
        const prevRepostedState = isReposted;
        const prevRepostsCount = repostsCount;

        try {
            setIsReposted(false);
            setRepostsCount(prev => Math.max(0, prev - 1));

            if (post.is_repost) {
                await repost(Number(post.parent_id));
            } else {
                await repost(Number(post.id));
            }
            setShowRepostOptions(false);
        } catch (error) {
            console.error('Error undoing repost:', error);
            setIsReposted(prevRepostedState);
            setRepostsCount(prevRepostsCount);
            alert('Failed to undo repost. Please try again.');
        } finally {
            setIsUndoingRepost(false);
        }
    };

    const openQuoteModal = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setShowQuoteModal(true);
        setShowRepostOptions(false);
    };

    const closeQuoteModal = () => {
        setShowQuoteModal(false);
    };

    const handleQuoteRepost = async (quoteText: string) => {
        if (isReposting || !quoteText.trim()) return;

        setIsReposting(true);
        const prevRepostedState = isReposted;
        const prevRepostsCount = repostsCount;

        try {
            setIsReposted(true);
            setRepostsCount(prev => prev + 1);

            if (post.is_repost) {
                await repost(Number(post.parent_id), quoteText);
            } else {
                await repost(Number(post.id), quoteText);
            }
            closeQuoteModal();
        } catch (error) {
            console.error('Error quote reposting:', error);
            setIsReposted(prevRepostedState);
            setRepostsCount(prevRepostsCount);
            alert('Failed to quote repost. Please try again.');
        } finally {
            setIsReposting(false);
        }
    };

    const openCommentModal = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setShowCommentModal(true);
    };

    const closeCommentModal = () => {
        setShowCommentModal(false);
    };

    if (loading) {
        return (
            <Navbar>
                <div className="flex items-center justify-center min-h-[200px]">
                    <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
                </div>
            </Navbar>
        );
    }

    if (error) {
        return (
            <Navbar>
                <div className="flex items-center justify-center min-h-[200px]">
                    <div className="text-center">
                        <h2 className="text-xl font-bold mb-2">Error</h2>
                        <p className="text-gray-400">{error}</p>
                    </div>
                </div>
            </Navbar>
        );
    }

    return (
        <Navbar>
            <div>
                <header className="sticky top-0 bg-black/80 backdrop-blur-sm p-4 flex items-center gap-6 border-b border-gray-800">
                    <button onClick={() => navigate(-1)} className="hover:bg-gray-800 p-2 rounded-full">
                        <ArrowLeft className="w-5 h-5" />
                    </button>
                    <h1 className="text-xl font-bold">Post</h1>
                </header>

                {/* Show parent post for comments */}
                {post.parent_post && !post.is_repost && (
                    <div className="border-b border-gray-800">
                        <Post
                            post={post.parent_post}
                            onRepost={async () => {
                                if (username && postId) {
                                    const updatedPost = await getPostById(postId);
                                    if (updatedPost) {
                                        setPost(updatedPost);
                                    }
                                }
                            }}
                        />
                    </div>
                )}

                <div className="p-2 flex items-start justify-between">
                    <div className="flex gap-3">
                        <div className="w-10 h-10 rounded-full bg-gray-700 flex-shrink-0 flex items-center justify-center text-lg font-bold">
                            {post.nickname?.charAt(0).toUpperCase() || post.username?.charAt(0).toUpperCase()}
                        </div>
                        <div>
                            <div className="flex items-center gap-1">
                                <span className="font-bold">{post.nickname || post.username}</span>
                            </div>
                            <div className="text-gray-500">@{post.username}</div>
                        </div>
                    </div>
                    <div className="flex gap-2">
                        {isCurrentUserAuthor && (
                            <div className="relative">
                                <button onClick={toggleOptions} className="p-2 hover:bg-gray-800 rounded-full">
                                    <MoreVertical className="w-5 h-5" />
                                </button>
                                {showOptions && (
                                    <PostOptions
                                        onEdit={handleEdit}
                                        onDelete={handleDelete}
                                        isDeleting={isDeleting}
                                    />
                                )}
                            </div>
                        )}
                    </div>
                </div>

                <div className="px-4 pb-4">
                    <PostContent
                        post={post}
                        isEditing={isEditing}
                        isQuoteRepost={isQuoteRepost}
                        isSimpleRepost={isSimpleRepost}
                        editContent={editContent}
                        handleContentChange={handleContentChange}
                        handleCancelEdit={handleCancelEdit}
                        handleSaveEdit={handleSaveEdit}
                        isSubmitting={isSubmitting}
                        error={editError}
                        MAX_CHARS={MAX_CHARS}
                        formatDate={formatDate}
                    />
                    <div className="text-gray-500 mt-4 ml-2">
                        <span>{formatDate(post.created_at)}</span>
                    </div>

                    <div className="flex gap-6 mt-4 px-3 py-4 border-y justify-between border-gray-800">
                        <button
                            className="flex items-center hover:text-blue-400"
                            onClick={openCommentModal}
                        >
                            <MessageCircle className="w-5 h-5 mr-1" />
                            <span>{commentsCount}</span>
                        </button>

                        <div className="relative">
                            <button
                                className={`flex items-center ${isReposted ? 'text-green-500' : 'hover:text-green-500'}`}
                                onClick={toggleRepostOptions}
                                disabled={isReposting || isUndoingRepost}
                            >
                                <Repeat
                                    size={18}
                                    className="mr-1"
                                    fill={isReposted ? 'currentColor' : 'none'}
                                />
                                <span>{repostsCount}</span>
                            </button>

                            {showRepostOptions && !isReposted && (
                                <div className="absolute bottom-full mb-2 right-0 bg-black rounded-lg shadow-lg border border-gray-800 z-10 overflow-hidden w-40">
                                    <button
                                        onClick={handleSimpleRepost}
                                        className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
                                    >
                                        <Repeat size={16} className="mr-2" />
                                        Repost
                                    </button>
                                    <button
                                        onClick={openQuoteModal}
                                        className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
                                    >
                                        <svg
                                            className="w-4 h-4 mr-2"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                            xmlns="http://www.w3.org/2000/svg"
                                        >
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth="2"
                                                d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"
                                            ></path>
                                        </svg>
                                        Quote
                                    </button>
                                </div>
                            )}
                        </div>

                        <button
                            className={`flex items-center ${isLiked ? 'text-red-500' : 'hover:text-red-500'}`}
                            onClick={handleLikeToggle}
                            disabled={isLikeProcessing}
                        >
                            <Heart
                                size={18}
                                className="mr-1"
                                fill={isLiked ? 'currentColor' : 'none'}
                            />
                            <span>{likesCount}</span>
                        </button>

                        <button className="flex items-center hover:text-blue-400">
                            <Share size={18} />
                        </button>
                    </div>

                    <div className="flex gap-3 pt-4">
                        <div className="w-10 h-10 rounded-full bg-gray-700 flex-shrink-0 flex items-center justify-center text-lg font-bold">
                            {currentUser?.nickname?.charAt(0).toUpperCase() || currentUser?.username?.charAt(0).toUpperCase() || 'U'}
                        </div>
                        <div className="flex-1 flex flex-col gap-2">
                            <textarea
                                value={newComment}
                                onChange={handleCommentChange}
                                placeholder="Post your reply"
                                className="w-full bg-transparent text-xl border-none outline-none text-white placeholder-gray-500 resize-none"
                                rows={2}
                            />
                            <div className="flex justify-between items-center">
                                <span className="text-gray-500 text-sm">
                                    {newComment.length}/{MAX_CHARS}
                                </span>
                                <button
                                    onClick={submitComment}
                                    disabled={!newComment.trim() || isAddingComment}
                                    className={`px-4 py-1.5 text-black rounded-full font-bold ${
                                        !newComment.trim() || isAddingComment
                                            ? 'bg-[#787A7A] text-black cursor-not-allowed'
                                            : 'bg-[#EFF3F4] hover:bg-[#D7DBDC] text-black'
                                    }`}
                                >
                                    {isAddingComment ? 'Posting...' : 'Reply'}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="border-b border-gray-800"></div>

                <div className="divide-y divide-gray-800">
                    {loadingComments ? (
                        <div className="flex justify-center py-8">
                            <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
                        </div>
                    ) : comments.length > 0 ? (
                        comments.map(comment => (
                            <Post
                                key={comment.id}
                                post={{
                                    ...comment,
                                    id: comment.id,
                                    created_at: comment.created_at,
                                    user_id: comment.user_id,
                                    is_repost: false,
                                    parent_post: comment.parent_post
                                }}
                            />
                        ))
                    ) : (
                        <div className="py-8 text-center text-gray-500">
                            No replies yet. Be the first to reply!
                        </div>
                    )}
                </div>

                {showCommentModal && (
                    <CommentModal
                        post={post}
                        currentUser={currentUser}
                        onClose={closeCommentModal}
                        onCommentAdded={handleCommentAdded}
                        formatDate={formatDate}
                        displayName={post.nickname || post.username}
                    />
                )}

                {showQuoteModal && (
                    <QuoteModal
                        post={post}
                        currentUser={currentUser}
                        onClose={closeQuoteModal}
                        onQuoteSubmit={handleQuoteRepost}
                        isReposting={isReposting}
                        formatDate={formatDate}
                        displayName={post.nickname || post.username}
                    />
                )}
            </div>
        </Navbar>
    );
}

export default PostDetailPage;