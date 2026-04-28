data Option a = None | Something a
foo :: a -> Option a
foo x = Something x