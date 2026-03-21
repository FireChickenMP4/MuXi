import Data.List (partition)   -- 引入 partition
import Data.Function ((&))      -- 引入 (&)，也可以不用

hylo :: Functor f => (f b -> b) -> (a -> f a) -> a -> b
hylo f g = f . fmap (hylo f g) . g

data BinTreeF a b = Tip | Branch b a b
    deriving (Functor)

quicksort :: Ord a => [a] -> [a]
quicksort = let
    split [] = Tip
    split (x:xs) = partition (< x) xs & \(l,r) -> Branch l x r
    -- or
    -- split (x:xs) = let (l, r) = partition (< x) xs in Branch l x r

    merge Tip = []
    merge (Branch l x r) = l ++ [x] ++ r
    in hylo merge split

main :: IO ()
main = do {
    print(quicksort [1,7,4,6,-2,0,114,514])
}