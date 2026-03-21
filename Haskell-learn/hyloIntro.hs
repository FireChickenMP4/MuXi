import Data.List (partition, unfoldr, sort)
import Data.Function ((&))      -- 引入 (&)，也可以不用
import System.Environment (getArgs)

-- 当然我们也可以实现类似c sort 的递归深度过深切换堆排序
-- 但是要在hylo中加入递归深度
hyloIntro :: Functor f => Int -> (f b -> b) -> (a -> f a) -> (a -> b) -> a -> b
hyloIntro 0 _ _ h x = h x
-- 深度耗尽改用备用排序
hyloIntro d f g h x = (f . fmap (hyloIntro (d-1) f g h) . g) x
-- hyloIntro d f g h x = f (fmap (hyloIntro (d-1) f g h) (g x))

data BinTreeF a b = Tip | Branch b a b
    deriving (Functor)

data LeftistHeap a = E | L a Int (LeftistHeap a) (LeftistHeap a)

rank :: LeftistHeap a -> Int
rank E = 0
rank (L _ r _ _) = r

merge :: Ord a => LeftistHeap a -> LeftistHeap a -> LeftistHeap a
merge h E = h
merge E h = h
merge h1@(L x r1 a1 b1) h2@(L y r2 a2 b2)
    | x <= y = makeT x a1 (merge b1 h2)
    | otherwise = makeT y a2 (merge b2 h1)
-- 大的放到右分支里面
-- makeT是一个辅助函数
-- 用于迭代rank
-- 以及维护s-value,在左偏堆里就是rankRight不超过rankLeft

makeT :: Ord a => a -> LeftistHeap a -> LeftistHeap a -> LeftistHeap a
makeT x a b = let (a', b') = if rank a >= rank b then (a,b) else (b,a)
              in L x (rank b' + 1) a' b'
              --这里右子树rank加一


insert :: Ord a => a -> LeftistHeap a -> LeftistHeap a
insert x h = merge (L x 1 E E) h
-- 这里直接交给合并了
-- 并初始化生成rank

findMin :: LeftistHeap a -> a
findMin (L x _ _ _) = x
findMin E = error "empty heap"
-- 取出来头节点就行了

deleteMin :: Ord a => LeftistHeap a -> LeftistHeap a
-- 这里deleteMin就是真delete了，而不是findMin + delete
deleteMin E = error "empty heap"
deleteMin (L _ _ a b) = merge a b

fromList :: Ord a => [a] -> LeftistHeap a
fromList = foldl (flip insert) E
-- foldl :: (b -> a -> b) -> b -> [a] -> b
-- 当前累计值b与列表中a一个元素结合，产生新的累计值
-- insert :: Ord a => a -> LeftistHeap a -> LeftistHeap a
-- flip insert :: Ord a => LeftistHeap a -> a -> LeftistHeap a

heapsort :: Ord a => [a] -> [a]
heapsort xs = unfoldr step (fromList xs)
    where
        step E = Nothing
        step h = Just (findMin h, deleteMin h)
-- 一个生成器函数step :: b -> Maybe (a, b)
-- 输入一个种子值b，可能返回Just (x, b') 产生一个元素和新的种子b'
-- 以及Nothing 表示生成结束
-- 然后反复调用 step，直到返回 Nothing，将所有产生的 x 依次放入列表中

quicksortIntro :: Ord a => [a] -> [a]
quicksortIntro xs = let
    maxDepth = floor (2 * logBase 2 (fromIntegral (length xs) + 1 :: Double))
    -- Int是i32 Integer类型是任意精度
    -- length xs 返回 Int
    -- 而 fromInteger 期望 Integer
    -- 应改用 fromIntegral 将 Int 转换为 Double（因为 logBase 需要浮点数）
    split [] = Tip
    split (x:xs) = partition (< x) xs & \(l,r) -> Branch l x r
    -- or
    -- split (x:xs) = let (l, r) = partition (< x) xs in Branch l x r

    merge Tip = []
    merge (Branch l x r) = l ++ [x] ++ r
    in hyloIntro maxDepth merge split heapsort xs

test :: IO ()
test = do
    let tests = [
            ([1,7,4,6,-2,0,114,514], sort [1,7,4,6,-2,0,114,514]),
            (replicate 1000 5, replicate 1000 5),
            ([1..10000],[1..10000]),
            ([10000,9999..1],[1..10000]),
            -- 负步长需要显式列出 否则得到空列表
            -- 第一个 第二个元素能让其推断步长
            ([10000000,9999999..1],[1..10000000])
            ] --这里必须整体括号缩进一下子
            --要不然会把runTest包进去，虽然没有，但是反正会报错
        runTest (input, expected) = do
                let result = quicksortIntro input
                putStr $ "Testing" ++ show (take 5 input) ++ "... "
                if result == expected
                    then putStrLn "Pass"
                    else do
                        putStrLn "Fail"
                        putStrLn $ "Got:        " ++ show (take 10 result) ++ "... "
                        putStrLn $ "Expected:   " ++ show (take 10 expected) ++ "... "
                -- 在 Haskell 中，$ 是一个函数应用运算符
                -- 它的主要作用是降低优先级并右结合
                -- 从而让我们可以省略括号，使代码更简洁易读。
    mapM_ runTest tests
        -- where -- do块里面不能用where
            -- 只用附加在顶层函数定义或者let表达式末尾
            -- 用于提供局部定义

main :: IO ()
main = do
    -- line <- getLine --读取一行
    -- -- 或者 contents <- getContents
    -- -- let nums = map read (lines contents) :: [Int]
    -- let nums = map read (words line) :: [Int]
    -- let sorted = quicksortIntro nums
    -- putStrLn $ unwords (map show sorted)
    -- words String -> [String]
    -- 将字符串按空白字符（空格、制表符、换行等）分割成单词列表，连续空白会被合并
    -- unwords [String] -> String
    -- 逆操作，单词列表连接成一个字符串，用空格隔开
    -- lines String -> [String]
    -- 将字符串按换行符 \n 分割成行列表
    -- unlines [String] -> String
    -- 同理，逆操作
    -- read Read a => String -> a
    -- 将字符串解析为指定类型的值
    -- 要求目标类型是 Read 类型类的实例（如 Int、Double、Bool 等）
    -- 由于 Haskell 的类型推断，通常需要上下文指定目标类型，例如 read "42" :: Int
    -- map read 将 read 函数映射到一个字符串列表上，得到对应类型的值列表
    -- map read ["1","2","3"] :: [Int] → [1,2,3]
    --在输入处理中，常与 words 结合：map read (words input) :: [Int] 
    -- 将空格分隔的字符串列表转换为整数列表
    -- mapM_ print xs 依次打印出xs中每个元素，每个占一行

    -- ineract (String -> String) -> IO ()
    -- 这是一个高阶函数，它将标准输入的所有内容作为字符串读入
    -- (需要用EOF，终端里是 <C-D>linux/mac <C-Z>win)
    -- 传递给参数函数（该函数处理输入字符串并返回输出字符串）
    -- 然后将输出字符串写入标准输出
    -- 它封装了 getContents 和 putStr 的常见模式
    -- 每行每列都可以
    -- interact $ \input ->
    --     let nums = map read (lines input) :: [Int]
    --         sorted = quicksortIntro nums
    --     in unlines (map show sorted) ++ "\n"
    -- interact $ \input ->
    --     let nums = map read (words input) :: [Int]
    --         sorted = quicksortIntro nums
    --     in unwords (map show sorted) ++ "\n"
    -- 但是说 words 已经能同时处理空格和换行分隔
    -- 所以没必要也
    -- 除非是要考虑空格的情况
    -- 类似getline(cin)

    -- 然后我要getArgs获取参数，如果是--test / -t 就调用test函数
    -- 否则接收输入
    args <- getArgs
    if "--test" `elem` args || "-t" `elem` args
        then test
        else do
            putStrLn "Space or Enter to separate"
            putStrLn "Ctrl+Z(Windows) / Ctrl+D(Linux Macos) To End Input"
            putStrLn "Input: "
            interact $ \input ->
                let nums = map read (words input) :: [Int]
                    sorted = quicksortIntro nums
                in unwords (map show sorted) ++ "\n"
    -- elem :: (Eq a) => a -> [a] -> bool
    -- 检查参数是否在列表中
    -- 而在Haskell中，任何普通双参数函数都可以通过``包围变成中缀运算符
    -- 更接近自然语言，可读性更好
    -- 有个notElem是相反的效果，检查是否不在