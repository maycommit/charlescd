---- MODULE circlerrRoutes ----
EXTENDS Naturals

CONSTANTS N

VARIABLES circles, pc, routes

ASSUME N \in Nat


NCircles == 1 .. N
Projects == {"P1", "P2", "P3"}


vars == << pc, circles, routes >>

Init ==
    /\ circles \in [NCircles -> SUBSET Projects]
    /\ pc \in [NCircles -> {"PGR", "HTY"}]
    /\ routes = [i \in NCircles |-> {}]

UponHealthy(self) ==
    /\ pc[self] = "HTY"
    /\ routes' = [routes EXCEPT ![self] = circles[self]]
    /\ circles' = circles
    /\ pc' = pc

Step(self) ==
    \/ UponHealthy(self)
    \/ UNCHANGED << pc, circles, routes >>

Next == (\E self \in NCircles: Step(self))

Spec == Init /\ [][Next]_vars
(*
    PGR -> Circle in progress
    DEG -> Circle error
    HTY -> Circle healthy 
*)

====