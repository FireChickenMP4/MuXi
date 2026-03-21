import Data.Typeable (typeOf)
instance Num (Int -> Int) where
    fromInteger n = \x -> x + (Prelude.fromInteger n :: Int)   -- 将整数 n 转换为“加 n”的函数
    (+) = undefined
    (-) = undefined
    (*) = undefined
    abs = undefined
    signum = undefined
    negate = undefined
class Appendable l a r where
    (#) :: l -> a -> [r]
instance Appendable () a a where
    () # x = [x]
instance Appendable [a] a a where
    xs # x = xs ++ [x]
instance Appendable [a -> b] a b where
    fs # x = [f x | f <- fs]
instance Appendable [a -> b] [a] b where
    fs # xs = [f x | f <- fs, x <-xs]
instance Appendable [()] a () where
    us # x = [() | _ <- us]
instance Appendable [()] [a] () where
    us # xs = [() | _ <- us, _ <- xs]
instance Appendable [[a]] a ([a],a) where
    lss # x = [(l,x) | l <- lss]
instance Appendable [[a]] [a] ([a],a) where
    lss # xs = [(l,x) | l <- lss,x <- xs]
fs = [(+1),(*2),(0 :: Int -> Int)]
units = [(),()]
lists = [[1],[17,19]]
xs = [2,3,5]
ys = [7,11,13]
main = do {
    print (xs >>= \x -> ys >>= \y -> return (x,y));
    print (fs <*> xs);
    print (units # xs);
    print (lists # xs)
}